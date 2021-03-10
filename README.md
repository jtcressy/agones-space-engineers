# agones-space-engineers
Running a dedicated space-engineers gameserver on Agones

## README ROUGH DRAFT - more edits and explanation to come. for now, this is an extremely basic overview of deployment steps.

# 1. Setup agones

# 2. Download dedicated server and initialize game world LOCALLY
- run program
- set console compatibility
- set password, server name, etc.
- start server then shut back down
- save as.. (name it SpaceEngineers-Dedicated.cfg and put it in this repo's working directory)
- Go to Saves directory take note of the folder that contains your world, use this in step 4!

# 3. `make all`
Build, push, deploy.

# 4. Copy world file to remote server
- `kubectl get pods -n gameservers` - get your pod's name
- `kubectl cp YourWorldSaveFolder <pod_name>:/appdata/space-engineers/SpaceEngineersDedicated/Saves`
- delete the pod to restart the server

# 5. Login with EOS network to your server!
You can browse the EOS server list on xbox or PC and you should see your server in the list.
