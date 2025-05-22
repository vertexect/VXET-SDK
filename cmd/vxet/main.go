// Copyright 2023 The VXET-Chain Authors
// This file is part of the VXET-Chain library.
//

// VXET-Client is a command-line client for VXET Chain.
package main

import (
	"fmt"
	"math/big"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/VXETChain/VXET-SDK/accounts"
	"github.com/VXETChain/VXET-SDK/accounts/keystore"
	"github.com/VXETChain/VXET-SDK/cmd/utils"
	"github.com/VXETChain/VXET-SDK/common"
	"github.com/VXETChain/VXET-SDK/console/prompt"
	"github.com/VXETChain/VXET-SDK/eth"
	"github.com/VXETChain/VXET-SDK/eth/downloader"
	"github.com/VXETChain/VXET-SDK/ethclient"
	"github.com/VXETChain/VXET-SDK/internal/debug"
	"github.com/VXETChain/VXET-SDK/internal/ethapi"
	"github.com/VXETChain/VXET-SDK/internal/flags"
	"github.com/VXETChain/VXET-SDK/log"
	"github.com/VXETChain/VXET-SDK/metrics"
	"github.com/VXETChain/VXET-SDK/node"
	"go.uber.org/automaxprocs/maxprocs"

	// Force-load the tracer engines to trigger registration
	_ "github.com/VXETChain/VXET-SDK/eth/tracers/js"
	_ "github.com/VXETChain/VXET-SDK/eth/tracers/native"

	"github.com/urfave/cli/v2"
)

const (
	clientIdentifier = "vxet" // Client identifier to advertise over the network
)

// Import from geth command
var (
	// Command definitions from geth/chaincmd.go, geth/accountcmd.go, etc.
	initCommand            = &cli.Command{ /*...*/ }
	importCommand          = &cli.Command{ /*...*/ }
	exportCommand          = &cli.Command{ /*...*/ }
	importHistoryCommand   = &cli.Command{ /*...*/ }
	exportHistoryCommand   = &cli.Command{ /*...*/ }
	importPreimagesCommand = &cli.Command{ /*...*/ }
	removedbCommand        = &cli.Command{ /*...*/ }
	dumpCommand            = &cli.Command{ /*...*/ }
	dumpGenesisCommand     = &cli.Command{ /*...*/ }
	accountCommand         = &cli.Command{ /*...*/ }
	walletCommand          = &cli.Command{ /*...*/ }
	consoleCommand         = &cli.Command{ /*...*/ }
	attachCommand          = &cli.Command{ /*...*/ }
	javascriptCommand      = &cli.Command{ /*...*/ }
	versionCommand         = &cli.Command{ /*...*/ }
	versionCheckCommand    = &cli.Command{ /*...*/ }
	licenseCommand         = &cli.Command{ /*...*/ }
	dumpConfigCommand      = &cli.Command{ /*...*/ }
	dbCommand              = &cli.Command{ /*...*/ }
	snapshotCommand        = &cli.Command{ /*...*/ }
	verkleCommand          = &cli.Command{ /*...*/ }
	logTestCommand         = &cli.Command{ /*...*/ }

	// Command line flags definition from geth/main.go
	nodeFlags    = []cli.Flag{ /*...*/ }
	rpcFlags     = []cli.Flag{ /*...*/ }
	consoleFlags = []cli.Flag{ /*...*/ }
	metricsFlags = []cli.Flag{ /*...*/ }
)

var app = flags.NewApp("VXET-Client Command Line Interface")

func init() {
	// Initialize the CLI app and start VXET-Client
	app.Action = vxetMain
	app.Commands = []*cli.Command{
		// Import all commands
		// ...omitted, should include all supported commands
	}
	if logTestCommand != nil {
		app.Commands = append(app.Commands, logTestCommand)
	}
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Flags = flags.Merge(
		nodeFlags,
		rpcFlags,
		consoleFlags,
		debug.Flags,
		metricsFlags,
	)
	flags.AutoEnvVars(app.Flags, "VXET")

	app.Before = func(ctx *cli.Context) error {
		maxprocs.Set() // Automatically set GOMAXPROCS to match Linux container CPU quota.
		flags.MigrateGlobalFlags(ctx)
		if err := debug.Setup(ctx); err != nil {
			return err
		}
		flags.CheckEnvVars(ctx, app.Flags, "VXET")
		return nil
	}
	app.After = func(ctx *cli.Context) error {
		debug.Exit()
		prompt.Stdin.Close() // Resets terminal mode.
		return nil
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// prepare manipulates memory cache allowance and setups metric system.
// This function should be called before launching devp2p stack.
func prepare(ctx *cli.Context) {
	// If we're running a known preset, log it for convenience.
	switch {
	case ctx.IsSet(utils.GoerliFlag.Name):
		log.Info("Starting VXET-Client on Görli testnet...")

	case ctx.IsSet(utils.SepoliaFlag.Name):
		log.Info("Starting VXET-Client on Sepolia testnet...")

	case ctx.IsSet(utils.HoleskyFlag.Name):
		log.Info("Starting VXET-Client on Holesky testnet...")

	case ctx.IsSet(utils.DeveloperFlag.Name):
		log.Info("Starting VXET-Client in ephemeral dev mode...")
		log.Warn(`You are running VXET-Client in --dev mode. Please note the following:

  1. This mode is only intended for fast, iterative development without assumptions on
     security or persistence.
  2. The database is created in memory unless specified otherwise. Therefore, shutting down
     your computer or losing power will wipe your entire block data and chain state for
     your dev environment.
  3. A random, pre-allocated developer account will be available and unlocked as
     eth.coinbase, which can be used for testing. The random dev account is temporary,
     stored on a ramdisk, and will be lost if your machine is restarted.
  4. Mining is enabled by default. However, the client will only seal blocks if transactions
     are pending in the mempool. The miner's minimum accepted gas price is 1.
  5. Networking is disabled; there is no listen-address, the maximum number of peers is set
     to 0, and discovery is disabled.
`)

	case !ctx.IsSet(utils.NetworkIdFlag.Name):
		log.Info("Starting VXET-Client on VXET mainnet...")
	}

	// Start metrics export if enabled
	utils.SetupMetrics(ctx)

	// Start system runtime metrics collection
	go metrics.CollectProcessMetrics(3 * time.Second)
}

// vxetMain is the main entry point into the system if no special subcommand is run.
// It creates a default node based on the command line arguments and runs it in
// blocking mode, waiting for it to be shut down.
func vxetMain(ctx *cli.Context) error {
	if args := ctx.Args().Slice(); len(args) > 0 {
		return fmt.Errorf("invalid command: %q", args[0])
	}

	prepare(ctx)
	stack, backend := makeFullNode(ctx)
	defer stack.Close()

	startNode(ctx, stack, backend, false)
	stack.Wait()
	return nil
}

// makeFullNode构建一个完整的eth节点服务
func makeFullNode(ctx *cli.Context) (*node.Node, ethapi.Backend) {
	stack, cfg := utils.MakeConfigNode(ctx)
	if ctx.IsSet(utils.OverrideCancun.Name) {
		cfg.Eth.OverrideCancun = new(big.Int).SetUint64(ctx.Uint64(utils.OverrideCancun.Name))
	}
	backend, eth := utils.RegisterEthService(stack, &cfg.Eth)

	// 配置基于不同的同步模式来注册
	if ctx.String(utils.SyncModeFlag.Name) == "light" {
		// 针对light客户端的特殊处理
		utils.RegisterLightFlags(ctx)
	}

	// 注册额外需要的服务
	utils.RegisterFilterAPI(stack, backend, eth)
	utils.RegisterShhService(stack)
	utils.RegisterLocalWallet(stack)
	utils.RegisterDNSDiscoveryService(stack, &cfg.Eth)

	return stack, backend
}

// startNode启动系统节点和所有注册的协议
func startNode(ctx *cli.Context, stack *node.Node, backend ethapi.Backend, isConsole bool) {
	debug.Memsize.Add("node", stack)

	// 启动节点本身
	utils.StartNode(ctx, stack, isConsole)

	// 解锁任何特别请求的账户
	unlockAccounts(ctx, stack)

	// 注册钱包事件处理程序
	events := make(chan accounts.WalletEvent, 16)
	stack.AccountManager().Subscribe(events)

	// 创建客户端以与本地节点交互
	rpcClient := stack.Attach()
	ethClient := ethclient.NewClient(rpcClient)

	go func() {
		// 打开已经附加的钱包
		for _, wallet := range stack.AccountManager().Wallets() {
			if err := wallet.Open(""); err != nil {
				log.Warn("无法打开钱包", "url", wallet.URL(), "err", err)
			}
		}
		// 监听钱包事件直到终止
		for event := range events {
			switch event.Kind {
			case accounts.WalletArrived:
				if err := event.Wallet.Open(""); err != nil {
					log.Warn("新钱包出现，无法打开", "url", event.Wallet.URL(), "err", err)
				}
			case accounts.WalletOpened:
				status, _ := event.Wallet.Status()
				log.Info("新钱包出现", "url", event.Wallet.URL(), "status", status)

				var derivationPaths []accounts.DerivationPath
				if event.Wallet.URL().Scheme == "ledger" {
					derivationPaths = append(derivationPaths, accounts.LegacyLedgerBaseDerivationPath)
				}
				derivationPaths = append(derivationPaths, accounts.DefaultBaseDerivationPath)

				event.Wallet.SelfDerive(derivationPaths, ethClient)

			case accounts.WalletDropped:
				log.Info("旧钱包丢弃", "url", event.Wallet.URL())
				event.Wallet.Close()
			}
		}
	}()

	// 为状态同步监控创建一个独立的goroutine，
	// 如果用户要求，同步完成后关闭节点。
	if ctx.Bool(utils.ExitWhenSyncedFlag.Name) {
		go func() {
			sub := stack.EventMux().Subscribe(downloader.DoneEvent{})
			defer sub.Unsubscribe()
			for {
				event := <-sub.Chan()
				if event == nil {
					continue
				}
				done, ok := event.Data.(downloader.DoneEvent)
				if !ok {
					continue
				}
				if timestamp := time.Unix(int64(done.Latest.Time), 0); time.Since(timestamp) < 10*time.Minute {
					log.Info("同步完成", "区块", done.Latest.Number, "哈希", done.Latest.Hash(),
						"时间", common.PrettyAge(timestamp))
					stack.Close()
				}
			}
		}()
	}

	// 如果启用，启动辅助服务
	if ctx.Bool(utils.MiningEnabledFlag.Name) {
		// 挖矿仅在运行完整的以太坊节点时有意义
		if ctx.String(utils.SyncModeFlag.Name) == "light" {
			utils.Fatalf("轻客户端不支持挖矿")
		}
		ethBackend, ok := backend.(*eth.EthAPIBackend)
		if !ok {
			utils.Fatalf("以太坊服务未运行")
		}
		// 设置来自CLI的Gas价格限制并开始挖矿
		gasprice := flags.GlobalBig(ctx, utils.MinerGasPriceFlag.Name)
		ethBackend.TxPool().SetGasTip(gasprice)
		if err := ethBackend.StartMining(); err != nil {
			utils.Fatalf("无法启动挖矿: %v", err)
		}
	}
}

// unlockAccounts解锁特别请求的账户
func unlockAccounts(ctx *cli.Context, stack *node.Node) {
	var unlocks []string
	inputs := strings.Split(ctx.String(utils.UnlockedAccountFlag.Name), ",")
	for _, input := range inputs {
		if trimmed := strings.TrimSpace(input); trimmed != "" {
			unlocks = append(unlocks, trimmed)
		}
	}
	// 如果没有账户需要解锁，则短路返回
	if len(unlocks) == 0 {
		return
	}
	// 如果节点的API暴露给外部，不允许不安全的账户解锁
	// 打印警告日志并跳过解锁
	if !stack.Config().InsecureUnlockAllowed && stack.Config().ExtRPCEnabled() {
		utils.Fatalf("不允许通过HTTP访问解锁账户!")
	}
	backends := stack.AccountManager().Backends(keystore.KeyStoreType)
	if len(backends) == 0 {
		log.Warn("无法解锁账户，keystore不可用")
		return
	}
	ks := backends[0].(*keystore.KeyStore)
	passwords := utils.MakePasswordList(ctx)
	for i, account := range unlocks {
		unlockAccount(ks, account, i, passwords)
	}
}

// unlockAccount解锁给定的账户
func unlockAccount(ks *keystore.KeyStore, address string, i int, passwords []string) {
	// 尝试使用提供的密码解锁
	if i < len(passwords) {
		if err := ks.Unlock(accounts.Account{Address: common.HexToAddress(address)}, passwords[i]); err == nil {
			log.Info("成功解锁账户", "address", address)
			return
		}
	}
	// 提示输入密码并持续尝试，直到成功或取消
	password := prompt.Stdin.PromptPassword("解锁账户 " + address + "的密码: ")
	if err := ks.Unlock(accounts.Account{Address: common.HexToAddress(address)}, password); err != nil {
		log.Error("无法解锁账户", "address", address, "err", err)
	} else {
		log.Info("成功解锁账户", "address", address)
	}
}
