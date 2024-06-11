package main

import (
	"os"
)

const ProfileIcon = "data:image/png;base64,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"

const ProfileType = "custom"

const ProfileName = "PodcrashPlay %s [%s]"

const JenkinsUrl = "https://jenkins.podcrash.com/job/PodcrashPlay/lastSuccessfulBuild/artifact/"

var MinecraftDir = map[string]string{
	"windows": ".minecraft",
	"darwin":  "minecraft",
}

const MinecraftVersionsDir = "versions"
const MinecraftLibrariesDir = "libraries"
const MinecraftLaunchwrapperDownload = "https://libraries.minecraft.net/net/minecraft/launchwrapper/1.12/launchwrapper-1.12.jar"
const MinecraftLaunchwrapperLocation = "net/minecraft/launchwrapper/1.12/"
const MinecraftLaunchwrapperName = "launchwrapper-1.12.jar"

const FileName = "%s-PodcrashPlay_b1.5.jar"

const PDCDownloadUrl = JenkinsUrl + FileName
const PDCDir = "PodcrashPlay"

const DefaultRAM = "2G"
const JavaArgs = "-Dlog4j2.formatMsgNoLookups=true -XX:HeapDumpPath=MojangTricksIntelDriversForPerformance_javaw.exe_minecraft.exe.heapdump -Xmx%s -XX:+UnlockExperimentalVMOptions -XX:+UseG1GC -XX:G1NewSizePercent=20 -XX:G1ReservePercent=20 -XX:MaxGCPauseMillis=50 -XX:G1HeapRegionSize=32M -Dpdc.basedir=\"%s\" -Dpdc.autoupdater=false -Dpdc.discord.id=708747697452089394 -javaagent:%s -XX:CompileCommand=exclude,java/lang/invoke/LambdaForm$*.invoke*"

const GameOptions = "Competitive:true\ngames.mineplex:true\ngames.minestrike:true\ngames.bridges:true\ngames.dominate:true\ngames.hypixel:true\ngames.copsvscrims:true"

const ForgeFile = "forge-%s-installer.jar"
const ForgeDownload = "https://maven.minecraftforge.net/net/minecraftforge/forge/%s/" + ForgeFile
const ForgeLibDir = "net/minecraftforge/forge/"

var ForgeVersion = map[string]string{
	"1.8.8": "1.8.8-11.15.0.1655",
	"1.8.9": "1.8.9-11.15.1.2318-1.8.9",
	"1.9":   "1.9-12.16.1.1887",
}

const OptifineFile = "OptiFine-" + OptifineLibraryVersionDir + ".jar"
const OptifineDownloadFile = "OptiFine_" + OptifineLibraryVersionDir + ".jar"
const OptifineVersionDir = "%s-OptiFine_HD_U_%s"
const OptifineLibraryVersionDir = "%s_HD_U_%s"
const OptifineUrl = "https://optifine.net/"
const OptifineDownload = OptifineUrl + "adloadx?f=" + OptifineDownloadFile
const OptifineLibraryDir = "optifine/OptiFine"
const OptifineLaunchwrapperDir = "optifine/launchwrapper-of/2.2"
const OptifineLaunchwrapperFile = "launchwrapper-of-2.2.jar"

var OptifineVersion = map[string]string{
	"1.8.8": "I7",
	"1.8.9": "M5",
	"1.9":   "I5",
}

var JreDownloads = map[string]string{
	"windows": "https://cdn.azul.com/zulu/bin/zulu8.78.0.19-ca-jre8.0.412-win_x64.zip",
	"darwin":  "https://cdn.azul.com/zulu/bin/zulu8.78.0.19-ca-jre8.0.412-macosx_x64.zip",
}

var JavaPath = map[string]string{
	"windows": "zulu8.78.0.19-ca-jre8.0.412-win_x64/bin/java.exe",
	"darwin":  "zulu8.78.0.19-ca-jre8.0.412-macosx_x64/zulu-8.jre/Contents/Home/bin/java",
}

var tempDir, _ = os.MkdirTemp("", "pdc")
var configDir, _ = os.UserConfigDir()
