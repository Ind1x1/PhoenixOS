package main

import (
	"fmt"
	"os"
	"time"

	"github.com/PhoenixOS-IPADS/PhOS/scripts/utils"
	"github.com/charmbracelet/log"
)

const (
	KLibClangPath    = "third_party/libclang-static-build"
	KLibYamlCppPath  = "third_party/yaml-cpp"
	KPhOSPath        = "pos"
	KPhOSCLIPath     = "pos/cli"
	KPhOSPatcherPath = "pos/cuda_impl/patcher"
	KRemotingPath    = "remoting/cuda"
	KBuildLogPath    = "build_log"
	KBuildLibPath    = "lib"
	KBuildIncPath    = "lib/pos/include"
	KBuildBinPath    = "bin"
	KInstallLibPath  = "/lib/x86_64-linux-gnu"
	KInstallIncPath  = "/usr/local/include"
	KInstallBinPath  = "/usr/local/bin"
)

func buildLibYamlCpp(cmdOpt CmdOptions, buildOpt BuildOptions, logger *log.Logger) {
	logger.Infof("building libyaml-cpp...")

	buildLogPath := fmt.Sprintf("%s/%s/%s", cmdOpt.RootDir, KBuildLogPath, "build_libyamlcpp.log")
	build_script := fmt.Sprintf(`
		#!/bin/bash
		set -e
		%s
		cd %s/%s
		rm -rf build
		mkdir build && cd build
		cmake -DYAML_BUILD_SHARED_LIBS=on .. >%s 2>&1
		make -j >%s 2>&1
		cp ./libyaml-cpp.so %s
		cp ./libyaml-cpp.so.0.8 %s
		cp ./libyaml-cpp.so.0.8.0 %s
		cp -r ../include/yaml-cpp %s/yaml-cpp
		`,
		buildOpt.export_string(),
		cmdOpt.RootDir, KLibYamlCppPath,
		buildLogPath,
		buildLogPath,
		fmt.Sprintf("%s/%s", cmdOpt.RootDir, KBuildLibPath),
		fmt.Sprintf("%s/%s", cmdOpt.RootDir, KBuildLibPath),
		fmt.Sprintf("%s/%s", cmdOpt.RootDir, KBuildLibPath),
		fmt.Sprintf("%s/%s", cmdOpt.RootDir, KBuildIncPath),
	)

	install_script := fmt.Sprintf(`
		#!/bin/bash
		set -e
		cd %s/%s
		cp ./libyaml-cpp.so %s/
		cp ./libyaml-cpp.so.0.8 %s/
		cp ./libyaml-cpp.so.0.8.0 %s/
		cp -r ./include %s
		`,
		cmdOpt.RootDir, KLibYamlCppPath,
		KInstallLibPath,
		KInstallLibPath,
		KInstallLibPath,
		KInstallIncPath,
	)

	start := time.Now()
	_, err := utils.BashScriptGetOutput(build_script, false, logger)
	if err != nil {
		logger.Fatalf("failed to build libyaml-cpp, please see log at %s", buildLogPath)
	}
	elapsed := time.Since(start)

	utils.ClearLastLine()
	logger.Infof("built libyaml-cpp: %.2fs", elapsed.Seconds())

	if *cmdOpt.DoInstall {
		_, err := utils.BashScriptGetOutput(install_script, false, logger)
		if err != nil {
			logger.Fatalf("failed to install libyaml-cpp, please see log at %s", buildLogPath)
		}
		logger.Infof("installed libyaml-cpp")
	}
}

func buildLibClang(cmdOpt CmdOptions, buildOpt BuildOptions, logger *log.Logger) {
	logger.Infof("building libclang...")

	buildLogPath := fmt.Sprintf("%s/%s/%s", cmdOpt.RootDir, KBuildLogPath, "build_libclang.log")
	build_script := fmt.Sprintf(`
		#!/bin/bash
		set -e
		%s
		cd %s/%s
		if [ ! -d "./build" ]; then
			mkdir build && cd build
			cmake .. -DCMAKE_INSTALL_PREFIX=.. >%s 2>&1
			make install -j >%s 2>&1
		fi
		cp ../lib/libclang.so %s
		cp ../lib/libclang.so.13 %s
		cp ../lib/libclang.so.VERSION %s
		cp -r ../include/clang-c %s/clang-c
		`,
		buildOpt.export_string(),
		cmdOpt.RootDir, KLibClangPath,
		buildLogPath,
		buildLogPath,
		fmt.Sprintf("%s/%s", cmdOpt.RootDir, KBuildLibPath),
		fmt.Sprintf("%s/%s", cmdOpt.RootDir, KBuildLibPath),
		fmt.Sprintf("%s/%s", cmdOpt.RootDir, KBuildLibPath),
		fmt.Sprintf("%s/%s", cmdOpt.RootDir, KBuildIncPath),
	)

	install_script := fmt.Sprintf(`
		#!/bin/bash
		set -e
		cd %s/%s
		cp ./lib/libclang.so %s/
		cp ./lib/libclang.so.13 %s/
		cp ./lib/libclang.so.VERSION %s/
		cp -r ./include/clang-c %s/clang-c
		`,
		cmdOpt.RootDir, KLibClangPath,
		KInstallLibPath,
		KInstallLibPath,
		KInstallLibPath,
		KInstallIncPath,
	)

	start := time.Now()
	_, err := utils.BashScriptGetOutput(build_script, false, logger)
	if err != nil {
		logger.Fatalf("failed to build libclang, please see log at %s", buildLogPath)
	}
	elapsed := time.Since(start)

	utils.ClearLastLine()
	logger.Infof("built libclang: %.2fs", elapsed.Seconds())

	if *cmdOpt.DoInstall {
		_, err := utils.BashScriptGetOutput(install_script, false, logger)
		if err != nil {
			logger.Fatalf("failed to install libclang, please see log at %s", buildLogPath)
		}
		logger.Infof("installed libclang")
	}
}

func buildKernelPatcher(cmdOpt CmdOptions, buildOpt BuildOptions, logger *log.Logger) {
	logger.Infof("building CUDA kernel patcher...")

	buildLogPath := fmt.Sprintf("%s/%s/%s", cmdOpt.RootDir, KBuildLogPath, "build_kernel_patcher.log")
	build_script := fmt.Sprintf(`
		#!/bin/bash
		set -e
		%s
		cd %s/%s
		if [ -d "./build" ]; then
			rm -rf build
		fi
		mkdir build && cd build
		cmake .. >%s 2>&1
		make -j >%s 2>&1
		if [ ! -e "./release/libpatcher.a" ] || [ ! -e "./patcher.h" ]; then
			exit 1
		fi
		cp ./release/libpatcher.a %s/%s
		cp ./patcher.h %s/%s
		`,
		buildOpt.export_string(),
		cmdOpt.RootDir, KPhOSPatcherPath,
		buildLogPath,
		buildLogPath,
		cmdOpt.RootDir, KBuildLibPath,
		cmdOpt.RootDir, KBuildIncPath,
	)

	install_script := fmt.Sprintf(`
		#!/bin/bash
		set -e
		cd %s/%s
		if [ -e "./release/libpatcher.a" ]; then
			cp ./release/libpatcher.a %s/
		else
			exit 1
		fi
		`,
		cmdOpt.RootDir, KLibClangPath,
		KInstallLibPath,
	)

	start := time.Now()
	_, err := utils.BashScriptGetOutput(build_script, false, logger)
	if err != nil {
		logger.Fatalf("failed to build CUDA kernel patcher, please see log at %s", buildLogPath)
	}
	elapsed := time.Since(start)

	utils.ClearLastLine()
	logger.Infof("built CUDA kernel patcher: %.2fs", elapsed.Seconds())

	if *cmdOpt.DoInstall {
		_, err := utils.BashScriptGetOutput(install_script, false, logger)
		if err != nil {
			logger.Fatalf("failed to install CUDA kernel patcher, please see log at %s", buildLogPath)
		}
		logger.Infof("installed CUDA kernel patcher")
	}
}

func buildPhOSCore(cmdOpt CmdOptions, buildOpt BuildOptions, logger *log.Logger) {
	logger.Infof("building PhOS core for CUDA target...")

	buildLogPath := fmt.Sprintf("%s/%s/%s", cmdOpt.RootDir, KBuildLogPath, "build_phos_core.log")
	build_script := fmt.Sprintf(`
		#!/bin/bash
		set -e
		%s
		cd %s
		rm -rf ./build
		# load build options
		meson build &>%s 2>&1
		cd build
		ninja clean
		ninja &>%s 2>&1
		cp %s/build/libpos.so %s/%s
		cp %s/build/pos/include/* %s/%s
		`,
		buildOpt.export_string(),
		cmdOpt.RootDir,
		buildLogPath,
		buildLogPath,
		cmdOpt.RootDir, cmdOpt.RootDir, KBuildLibPath,
		cmdOpt.RootDir, cmdOpt.RootDir, KBuildIncPath,
	)

	install_script := fmt.Sprintf(`
		#!/bin/bash
		set -e
		cp %s/build/*.so %s
		`,
		cmdOpt.RootDir, KInstallLibPath,
	)

	start := time.Now()
	_, err := utils.BashScriptGetOutput(build_script, false, logger)
	if err != nil {
		logger.Fatalf("failed to build PhOS Core for CUDA target, please see log at %s", buildLogPath)
	}
	elapsed := time.Since(start)

	utils.ClearLastLine()
	logger.Infof("built PhOS Core for CUDA target: %.2fs", elapsed.Seconds())

	if *cmdOpt.DoInstall {
		_, err := utils.BashScriptGetOutput(install_script, false, logger)
		if err != nil {
			logger.Fatalf("failed to install PhOS core, please see log at %s", buildLogPath)
		}
		logger.Infof("installed PhOS core")
	}
}

func buildPhOSCLI(cmdOpt CmdOptions, buildOpt BuildOptions, logger *log.Logger) {
	logger.Infof("building PhOS CLI...")

	buildLogPath := fmt.Sprintf("%s/%s/%s", cmdOpt.RootDir, KBuildLogPath, "build_phos_cli.log")
	build_script := fmt.Sprintf(`
		#!/bin/bash
		set -e
		%s
		cd %s/%s
		rm -rf build
		mkdir build
		cd build
		cmake .. &>%s 2>&1
		make -j  &>%s 2>&1
		cp ./pos-cli %s/%s
		`,
		buildOpt.export_string(),
		cmdOpt.RootDir, KPhOSCLIPath,
		buildLogPath,
		buildLogPath,
		cmdOpt.RootDir, KBuildBinPath,
	)

	install_script := fmt.Sprintf(`
		#!/bin/bash
		set -e
		cd %s/%s
		cp ./build/pos-cli %s
		`,
		cmdOpt.RootDir, KPhOSCLIPath,
		KInstallBinPath,
	)

	start := time.Now()
	_, err := utils.BashScriptGetOutput(build_script, false, logger)
	if err != nil {
		logger.Fatalf("failed to build PhOS CLI: %s", err)
	}
	elapsed := time.Since(start)

	utils.ClearLastLine()
	logger.Infof("built PhOS CLI: %.2fs", elapsed.Seconds())

	if *cmdOpt.DoInstall {
		_, err := utils.BashScriptGetOutput(install_script, false, logger)
		if err != nil {
			logger.Fatalf("failed to install PhOS CLI, please see log at %s", buildLogPath)
		}
		logger.Infof("installed PhOS CLI")
	}
}

func buildRemoting(cmdOpt CmdOptions, buildOpt BuildOptions, logger *log.Logger) {
	logger.Infof("building remoting framework...")

	buildLogPath := fmt.Sprintf("%s/%s/%s", cmdOpt.RootDir, KBuildLogPath, "build_remoting_framework.log")
	build_script := fmt.Sprintf(`
		#!/bin/bash
		set -e
		%s
		cd %s/%s
		make libtirpc -j &>%s 2>&1
		cp ./submodules/libtirpc/install/lib/libtirpc.so %s/%s
		cd cpu
		make clean
		LOG=DEBUG make cricket-rpc-server cricket-client.so -j &>%s 2>&1
		cp cricket-rpc-server %s/%s
		cp cricket-client.so %s/%s
		`,
		buildOpt.export_string(),
		cmdOpt.RootDir, KRemotingPath,
		buildLogPath,
		cmdOpt.RootDir, KBuildLibPath,
		buildLogPath,
		cmdOpt.RootDir, KBuildBinPath,
		cmdOpt.RootDir, KBuildLibPath,
	)

	install_script := fmt.Sprintf(`
		#!/bin/bash
		set -e
		cd %s/%s/cpu
		cp cricket-rpc-server %s
		cp cricket-client.so %s
		`,
		cmdOpt.RootDir, KRemotingPath,
		KInstallBinPath,
		KInstallLibPath,
	)

	start := time.Now()
	_, err := utils.BashScriptGetOutput(build_script, false, logger)
	if err != nil {
		logger.Fatalf("failed to build remoting framework: %s", err)
	}
	elapsed := time.Since(start)

	utils.ClearLastLine()
	logger.Infof("built remoting framework: %.2fs", elapsed.Seconds())

	if *cmdOpt.DoInstall {
		_, err := utils.BashScriptGetOutput(install_script, false, logger)
		if err != nil {
			logger.Fatalf("failed to install remoting framework, please see log at %s", buildLogPath)
		}
		logger.Infof("installed remoting framework")
	}
}

func buildUnitTest(cmdOpt CmdOptions, buildOpt BuildOptions, logger *log.Logger) {
	// logger.Infof("building uint test...")
}

func cleanCommon(cmdOpt CmdOptions, logger *log.Logger) {
	logger.Infof("cleaning common directoroies...")
	clean_script := ""
	if *cmdOpt.WithThirdParty {
		clean_script = fmt.Sprintf(`
			#!/bin/bash
			rm -rf %s/%s/*
			rm -rf %s/%s/*
			rm -rf %s/%s/*
			rm -rf %s/%s/*
			`,
			cmdOpt.RootDir, KBuildBinPath,
			cmdOpt.RootDir, KBuildLibPath,
			cmdOpt.RootDir, KBuildLogPath,
			cmdOpt.RootDir, KBuildIncPath,
		)
	} else {
		clean_script = fmt.Sprintf(`
			#!/bin/bash
			rm -rf %s/%s/*
			rm -rf %s/%s/libpos.so
			# TODO: just for fast compilation of PhOS, decomment later
			# rm -rf %s/%s/libpatcher.a
			# rm -rf %s/%s/*.h
			rm -rf %s/%s/*
			`,
			cmdOpt.RootDir, KBuildBinPath,
			cmdOpt.RootDir, KBuildLibPath,
			cmdOpt.RootDir, KBuildLibPath,
			cmdOpt.RootDir, KBuildIncPath,
			cmdOpt.RootDir, KBuildLogPath,
		)
	}

	_, err := utils.BashScriptGetOutput(clean_script, true, logger)
	if err != nil {
		logger.Warnf("failed to clean common directoroies")
	} else {
		logger.Infof("done")
	}
}

func cleanLibYamlCpp(cmdOpt CmdOptions, logger *log.Logger) {
	logger.Infof("cleaning libyaml-cpp...")
	clean_script := fmt.Sprintf(`
		#!/bin/bash
		cd %s/%s
		rm -rf build
		rm -f %s/libyaml-cpp.so
		rm -f %s/libyaml-cpp.so.0.8
		rm -f %s/libyaml-cpp.so.0.8.0
		rm -rf %s/yaml-cpp
		`,
		cmdOpt.RootDir, KLibYamlCppPath,
		KInstallLibPath,
		KInstallLibPath,
		KInstallLibPath,
		KInstallIncPath,
	)
	_, err := utils.BashScriptGetOutput(clean_script, true, logger)
	if err != nil {
		logger.Warnf("failed to clean libyaml-cpp")
	}
	logger.Infof("done")
}

func cleanLibClang(cmdOpt CmdOptions, logger *log.Logger) {
	logger.Infof("cleaning libclang...")
	clean_script := fmt.Sprintf(`
		#!/bin/bash
		cd %s/%s
		rm -rf build
		rm -rf include
		rm -rf lib
		rm -rf share
		rm -f %s/libclang.so
		rm -f %s/libclang.so.13
		rm -f %s/libclang.so.VERSION
		rm -rf %s/clang-c
		`,
		cmdOpt.RootDir, KLibClangPath,
		KInstallLibPath,
		KInstallLibPath,
		KInstallLibPath,
		KInstallIncPath,
	)
	_, err := utils.BashScriptGetOutput(clean_script, true, logger)
	if err != nil {
		logger.Warnf("failed to clean libclang")
	}
	logger.Infof("done")
}

func cleanKernelPatcher(cmdOpt CmdOptions, logger *log.Logger) {
	logger.Infof("cleaning CUDA kernel patcher...")
	clean_script := fmt.Sprintf(`
		#!/bin/bash
		cd %s/%s
		rm -rf build
		rm -f %s/patcher.h
		rm -f %s/libpatcher.a
		`,
		cmdOpt.RootDir, KPhOSPatcherPath,
		KInstallIncPath,
		KInstallLibPath,
	)
	_, err := utils.BashScriptGetOutput(clean_script, true, logger)
	if err != nil {
		logger.Warnf("failed to clean libclang")
	}
	logger.Infof("done")
}

func cleanPhOSCore(cmdOpt CmdOptions, logger *log.Logger) {
	logger.Infof("cleaning PhOS core...")
	clean_script := fmt.Sprintf(`
		#!/bin/bash
		cd %s
		rm -rf build
		rm -f %s/libpos.so
		`,
		cmdOpt.RootDir,
		KInstallLibPath,
	)
	_, err := utils.BashScriptGetOutput(clean_script, true, logger)
	if err != nil {
		logger.Fatalf("failed to clean PhOS core")
	}
	logger.Infof("done")
}

func cleanPhOSCLI(cmdOpt CmdOptions, logger *log.Logger) {
	logger.Infof("cleaning PhOS CLI...")
	clean_script := fmt.Sprintf(`
		#!/bin/bash
		cd %s/%s
		rm -rf build
		rm -f %s/pos-cli
		`,
		cmdOpt.RootDir, KPhOSCLIPath,
		KInstallBinPath,
	)
	_, err := utils.BashScriptGetOutput(clean_script, true, logger)
	if err != nil {
		logger.Fatalf("failed to clean PhOS CLI")
	}
	logger.Infof("done")
}

func BuildTarget_CUDA(cmdOpt CmdOptions, buildOpt BuildOptions, logger *log.Logger) {
	// ==================== Prepare ====================
	logger.Infof("pre-build check...")
	utils.CheckAndInstallCommand("git", "git", nil, logger)
	utils.CheckAndInstallCommand("gcc", "build-essential", nil, logger)
	utils.CheckAndInstallCommand("g++", "build-essential", nil, logger)
	utils.CheckAndInstallCommand("yes", "yes", nil, logger)
	utils.CheckAndInstallCommand("cmake", "cmake", nil, logger)
	utils.CheckAndInstallCommand("curl", "curl", nil, logger)
	utils.CheckAndInstallCommand("tar", "tar", nil, logger)
	utils.CheckAndInstallCommand("tmux", "tmux", nil, logger)

	install_meson := func() error {
		_, err := utils.BashScriptGetOutput(`
			#!/bin/bash
			set -e
			pip3 install meson
			`, false, logger,
		)
		return err
	}
	utils.CheckAndInstallCommand("meson", "", install_meson, logger)

	install_ninja := func() error {
		_, err := utils.BashScriptGetOutput(`
			#!/bin/bash
			set -e
			pip3 install ninja
			`, false, logger,
		)
		return err
	}
	utils.CheckAndInstallCommand("ninja", "", install_ninja, logger)

	build_cargo := func() error {
		_, err := utils.BashScriptGetOutput(`
			#!/bin/bash
			set -e
			if tmux has-session -t cargo_installer 2>/dev/null; then
				tmux kill-session -t cargo_installer
			fi
			tmux new -s cargo_installer -d
			tmux send -t cargo_installer "curl https://sh.rustup.rs -sSf | sh; exit 0" ENTER
			tmux send-keys -t cargo_installer C-m
			echo '. "$HOME/.cargo/env"' >> $HOME/.bashrc
			`,
			false, logger,
		)
		return err
	}
	utils.CheckAndInstallCommand("cargo", "", build_cargo, logger)

	buildLogPath := fmt.Sprintf("%s/%s", cmdOpt.RootDir, KBuildLogPath)
	if err := utils.CreateDir(buildLogPath, false, 0775, logger); err != nil && !os.IsExist(err) {
		logger.Fatalf("failed to create directory for build logs at %s", buildLogPath)
	}

	libPath := fmt.Sprintf("%s/%s", cmdOpt.RootDir, KBuildLibPath)
	if err := utils.CreateDir(libPath, false, 0775, logger); err != nil && !os.IsExist(err) {
		logger.Fatalf("failed to create directory for built lib at %s", libPath)
	}

	includePath := fmt.Sprintf("%s/%s", cmdOpt.RootDir, KBuildIncPath)
	if err := utils.CreateDir(includePath, false, 0775, logger); err != nil && !os.IsExist(err) {
		logger.Fatalf("failed to create directory for built headers at %s", includePath)
	}

	binPath := fmt.Sprintf("%s/%s", cmdOpt.RootDir, KBuildBinPath)
	if err := utils.CreateDir(binPath, false, 0775, logger); err != nil && !os.IsExist(err) {
		logger.Fatalf("failed to create directory for built binary at %s", binPath)
	}

	// ==================== Build Dependencies ====================
	if *cmdOpt.WithThirdParty {
		logger.Infof("building dependencies...")
		buildLibClang(cmdOpt, buildOpt, logger)
		buildLibYamlCpp(cmdOpt, buildOpt, logger)

		// TODO: just for fast compilation of PhOS, remove later
		buildKernelPatcher(cmdOpt, buildOpt, logger)
	}

	// ==================== Build PhOS ====================
	// buildKernelPatcher(cmdOpt, logger)
	buildPhOSCore(cmdOpt, buildOpt, logger)
	buildPhOSCLI(cmdOpt, buildOpt, logger)
	buildRemoting(cmdOpt, buildOpt, logger)

	// ==================== Build and Run Unit Test ====================
	if *cmdOpt.DoUnitTest {
		BuildGoogleTest(cmdOpt, buildOpt, logger)
		buildUnitTest(cmdOpt, buildOpt, logger)
	}
}

func CleanTarget_CUDA(cmdOpt CmdOptions, logger *log.Logger) {
	// ==================== Clean Dependencies ====================
	if *cmdOpt.WithThirdParty {
		logger.Infof("cleaning dependencies...")
		cleanLibClang(cmdOpt, logger)
		cleanLibYamlCpp(cmdOpt, logger)

		// TODO: just for fast compilation of PhOS, remove later
		cleanKernelPatcher(cmdOpt, logger)
	}

	// ==================== Clean PhOS ====================
	cleanCommon(cmdOpt, logger)
	// cleanKernelPatcher(cmdOpt, logger)
	cleanPhOSCore(cmdOpt, logger)
	cleanPhOSCLI(cmdOpt, logger)

	// ==================== Clean Unit Test ====================
	CleanGoogleTest(cmdOpt, logger)
}
