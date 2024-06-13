REM #pack.bat
@echo off
setlocal enabledelayedexpansion

@REM set "names=理赔内网管理-beta 理赔管理-beta 理赔内网录入-beta 理赔录入-beta"
set "names=理赔内网录入-beta  理赔录入-beta"
@REM set "names=理赔内网管理-beta 理赔管理-beta "

for %%i in (%names%) do (

    ECHO 开始替换package.json -- %%i
    (for /f "tokens=*" %%a in (package.json.tmp) do (
        set line=%%a
        set line=!line:[Name]=%%i!
        echo !line!
    )) > package.json
    ECHO 替换完成package.json -- %%i

    ECHO 开始替换config.js -- %%i
    (for /f "tokens=*" %%a in (config.js.tmp) do (
        set line=%%a
        set line=!line:[Name]=%%i!
        echo !line!
    )) > config.js
    ECHO 替换完成config.js -- %%i

    ECHO 开始打包 ------------------------------------------------ %%i
    if %1==32 (
        ECHO 32位系统
        call npm run package -- --arch="ia32"
    ) else (
        ECHO 64位系统
        call npm run package
    )
    ECHO 打包完成 ------------------------------------------------ %%i
    @REM del package.json
    @REM del config.js
)

endlocal
@REM pause>nul
@REM exit