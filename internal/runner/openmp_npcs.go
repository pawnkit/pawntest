package runner

import (
	"maps"

	"github.com/pawnkit/pawntest/internal/backend"
)

type testNPC struct {
	name                            string
	spawned, dead, invulnerable     bool
	moving, movingToPlayer          bool
	position, rotation, velocity    [3]float32
	moveTarget                      [3]float32
	facing, health, armour          float32
	world, skin, interior           int
	movePlayer                      int
	vehicle, seat                   int
	weapon, ammo, ammoInClip        int
	keys                            [3]int
	melee, reloadEnabled            bool
	reloading, infiniteAmmo         bool
	weaponState, fightingStyle      int
	shooting, aiming                bool
	aimPlayer                       int
	aimPoint                        [3]float32
	weaponAccuracy                  map[int]float32
	weaponReloadTime                map[int]int
	weaponShootTime                 map[int]int
	weaponClipSize                  map[int]int
	weaponSkills                    map[int]int
	animationID                     int
	animation                       actorAnimation
	hasAnimation                    bool
	specialAction                   int
	playbackName                    string
	playbackRecord                  int
	playingPlayback, pausedPlayback bool
	surfingOffset                   [3]float32
	surfingVehicle, surfingObject   int
	surfingPlayerObject             int
	currentPath, currentPathPoint   int
	playingNode, nodeID             int
	nodePaused                      bool
}

type npcPathPoint struct {
	position  [3]float32
	stopRange float32
}

type npcPath struct{ points []npcPathPoint }

type npcNode struct {
	open  bool
	point int
}

type npcState struct {
	next       int
	npcs       map[int]*testNPC
	players    *openMPState
	vehicles   *vehicleState
	nextRecord int
	records    map[int]string
	nextPath   int
	paths      map[int]*npcPath
	nodes      map[int]*npcNode
}

func newNPCState() *npcState {
	return &npcState{npcs: map[int]*testNPC{}, records: map[int]string{}, paths: map[int]*npcPath{}, nodes: map[int]*npcNode{}}
}

func (state *npcState) Clone() scenarioModule {
	clone := newNPCState()

	clone.next = state.next
	clone.nextRecord, clone.nextPath = state.nextRecord, state.nextPath
	maps.Copy(clone.records, state.records)

	for id, path := range state.paths {
		clone.paths[id] = &npcPath{points: append([]npcPathPoint(nil), path.points...)}
	}

	for id, node := range state.nodes {
		nodeCopy := *node
		clone.nodes[id] = &nodeCopy
	}

	for id, npc := range state.npcs {
		npcCopy := *npc
		npcCopy.weaponAccuracy = maps.Clone(npc.weaponAccuracy)
		npcCopy.weaponReloadTime = maps.Clone(npc.weaponReloadTime)
		npcCopy.weaponShootTime = maps.Clone(npc.weaponShootTime)
		npcCopy.weaponClipSize = maps.Clone(npc.weaponClipSize)
		npcCopy.weaponSkills = maps.Clone(npc.weaponSkills)
		clone.npcs[id] = &npcCopy
	}

	return clone
}

func (state *npcState) Register(vm backend.VM, context *executionContext) error {
	state.players = context.scenarios.playerState()
	state.vehicles = context.scenarios.vehicleState()

	return registerScenarioNatives(vm, state.natives(context.state), context.mocks, context.allowUnknown)
}

func (state *npcState) natives(result *nativeState) map[string]backend.NativeFunc {
	natives := map[string]backend.NativeFunc{
		"__pt_npc_create": state.createNPC, "__pt_npc_valid": state.assertValid(result),
		"__pt_npc_spawned": state.assertSpawned(result), "__pt_npc_health": state.assertHealth(result),
		"__pt_npc_pos_near": state.assertPosition(result), "NPC_Create": state.createNPC,
		"NPC_Destroy": state.destroyNPC, "NPC_IsValid": state.isValidNPC,
		"NPC_IsDead": state.isDeadNPC, "NPC_Spawn": state.spawnNPC,
		"NPC_Respawn": state.spawnNPC, "NPC_IsSpawned": state.isSpawnedNPC,
		"NPC_GetAll": state.getAllNPCs, "NPC_SetPos": state.setNPCVector(npcPosition),
		"NPC_GetPos": state.getNPCVector(npcPosition), "NPC_SetRot": state.setNPCVector(npcRotation),
		"NPC_GetRot": state.getNPCVector(npcRotation), "NPC_SetVelocity": state.setNPCVector(npcVelocity),
		"NPC_GetVelocity": state.getNPCVector(npcVelocity), "NPC_GetPosMovingTo": state.getNPCVector(npcMoveTarget),
		"NPC_SetFacingAngle": state.setNPCFloat(npcFacing), "NPC_GetFacingAngle": state.getNPCFloatOutput(npcFacing),
		"NPC_SetVirtualWorld": state.setNPCInt(npcWorld), "NPC_GetVirtualWorld": state.getNPCInt(npcWorld),
		"NPC_SetSkin": state.setNPCInt(npcSkin), "NPC_GetSkin": state.getNPCInt(npcSkin),
		"NPC_GetCustomSkin": state.getNPCInt(npcSkin), "NPC_SetInterior": state.setNPCInt(npcInterior),
		"NPC_GetInterior": state.getNPCInt(npcInterior), "NPC_SetHealth": state.setNPCFloat(npcHealth),
		"NPC_GetHealth": state.getNPCFloat(npcHealth), "NPC_SetArmour": state.setNPCFloat(npcArmour),
		"NPC_GetArmour": state.getNPCFloat(npcArmour), "NPC_SetInvulnerable": state.setNPCInvulnerable,
		"NPC_IsInvulnerable": state.isNPCInvulnerable, "NPC_Kill": state.killNPC,
		"NPC_IsStreamedIn": state.isNPCStreamedIn, "NPC_IsAnyStreamedIn": state.isNPCAnyStreamedIn,
		"NPC_Move": state.moveNPC, "NPC_MoveToPlayer": state.moveNPCToPlayer,
		"NPC_StopMove": state.stopNPCMove, "NPC_IsMoving": state.isNPCMoving,
		"NPC_IsMovingToPlayer": state.isNPCMovingToPlayer, "NPC_SetAngleToPos": state.setNPCAngleToPosition,
		"NPC_SetAngleToPlayer": state.setNPCAngleToPlayer, "NPC_PutInVehicle": state.putNPCInVehicle,
		"NPC_RemoveFromVehicle": state.removeNPCFromVehicle, "NPC_GetVehicle": state.getNPCVehicle,
		"NPC_GetVehicleID": state.getNPCVehicle, "NPC_GetVehicleSeat": state.getNPCVehicleSeat,
		"NPC_EnterVehicle": state.putNPCInVehicle, "NPC_ExitVehicle": state.removeNPCFromVehicle,
	}
	maps.Copy(natives, state.combatNatives(result))
	maps.Copy(natives, state.navigationNatives(result))

	return natives
}

func (state *npcState) createNPC(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	if len(params) == 0 {
		return -1, nil
	}

	name, err := ctx.ReadString(params[0])
	if err != nil {
		return -1, err
	}

	id := state.next
	state.next++
	state.npcs[id] = &testNPC{
		name: name, health: 100, movePlayer: -1, vehicle: -1, seat: -1, aimPlayer: -1,
		playbackRecord: -1, surfingVehicle: -1, surfingObject: -1, surfingPlayerObject: -1,
		currentPath: -1, currentPathPoint: -1, playingNode: -1, nodeID: -1,
		weaponAccuracy: map[int]float32{}, weaponReloadTime: map[int]int{}, weaponShootTime: map[int]int{},
		weaponClipSize: map[int]int{}, weaponSkills: map[int]int{},
	}

	return backend.Cell(id), nil
}

func (state *npcState) npc(params []backend.Cell) (*testNPC, bool) {
	if len(params) == 0 {
		return nil, false
	}

	npc, ok := state.npcs[int(params[0])]

	return npc, ok
}

func npcPosition(npc *testNPC) *[3]float32   { return &npc.position }
func npcRotation(npc *testNPC) *[3]float32   { return &npc.rotation }
func npcVelocity(npc *testNPC) *[3]float32   { return &npc.velocity }
func npcMoveTarget(npc *testNPC) *[3]float32 { return &npc.moveTarget }
func npcFacing(npc *testNPC) *float32        { return &npc.facing }
func npcHealth(npc *testNPC) *float32        { return &npc.health }
func npcArmour(npc *testNPC) *float32        { return &npc.armour }
func npcWorld(npc *testNPC) *int             { return &npc.world }
func npcSkin(npc *testNPC) *int              { return &npc.skin }
func npcInterior(npc *testNPC) *int          { return &npc.interior }
