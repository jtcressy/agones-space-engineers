#!/bin/bash
# VARIABLES
GAME_DIR="/appdata/space-engineers/SpaceEngineersDedicated"
CONFIG_PATH="/${GAME_DIR}/SpaceEngineers-Dedicated.cfg"
INSTANCE_IP="0.0.0.0"

echo "-------------------------------START------------------------------"
/usr/games/steamcmd +login anonymous +@sSteamCmdForcePlatformType windows +force_install_dir ${GAME_DIR} +app_update 298740 +quit

cd ${GAME_DIR}/DedicatedServer64/
wine SpaceEngineersDedicated.exe -noconsole -ignorelastsession -path Z:\\appdata\\space-engineers\\SpaceEngineersDedicated
echo "--------------------------------END-------------------------------"
sleep 10