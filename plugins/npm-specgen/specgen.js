#! /usr/bin/env node

const path = require('path');
const fs = require('fs');
const child_process = require("child_process");
const util = require('util');
//const exec = util.promisify(child_process.exec);
const exec = child_process.exec

const getOsName = () => {
    switch(process.platform) {
        case `darwin`:
        case `linux`:
            return process.platform
        case `win32`:
            return `windows`
        default:
            throw Error(`Unsupported platform: ${process.platform}`)
    }
}

const getExeName = (toolName) => {
    if (getOsName() === "windows") {
        return `${toolName}.exe`
    }
    return toolName
}

const getArch = () => {
    if (process.arch === "x64")
        return "amd64"
    throw Error(`Unsupported architecture: ${process.arch}`)
}

const getSpecgenPath = () => {
    const osname = getOsName()
    const arch = getArch()

    const specgenPath = path.join(__dirname, `/dist/${osname}_${arch}/${getExeName("specgen")}`)
    if (!fs.existsSync(specgenPath)) {
        throw Error(`Can't find specgen tool at ${specgenPath}`)
    }
    return specgenPath
}

const runSpecgen = (specgenCommand) => {
    const specgenCommandLine = specgenCommand.join(" ")
    console.log(`Running specgen tool: ${specgenCommandLine}`)

    exec(specgenCommandLine, (error, stdout, stderr) => {
        console.error(stderr);
        console.log(stdout);
        if (error) {
            console.error(`Specgen tool raised error, exit code: ${error.code}`);
            console.error(error.message);
            if (error.code !== 0) { process.exit(error.code) }
        }
    })
}

var args = process.argv.slice(2);
const specgenPath = getSpecgenPath()
args.unshift(specgenPath)
runSpecgen(args)