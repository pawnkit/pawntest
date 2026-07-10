package runner

import "github.com/pawnkit/pawntest/internal/backend"

func (state *npcState) combatNatives(result *nativeState) map[string]backend.NativeFunc {
	return map[string]backend.NativeFunc{
		"__pt_npc_weapon": state.assertWeapon(result), "__pt_npc_aiming": state.assertAiming(result),
		"__pt_npc_animation": state.assertAnimation(result), "NPC_SetWeapon": state.setNPCInt(npcWeapon),
		"NPC_GetWeapon": state.getNPCInt(npcWeapon), "NPC_SetAmmo": state.setNPCInt(npcAmmo),
		"NPC_GetAmmo": state.getNPCInt(npcAmmo), "NPC_SetAmmoInClip": state.setNPCInt(npcAmmoInClip),
		"NPC_GetAmmoInClip": state.getNPCInt(npcAmmoInClip), "NPC_SetFightingStyle": state.setNPCInt(npcFightingStyle),
		"NPC_GetFightingStyle": state.getNPCInt(npcFightingStyle), "NPC_SetWeaponState": state.setNPCInt(npcWeaponState),
		"NPC_GetWeaponState": state.getNPCInt(npcWeaponState), "NPC_SetSpecialAction": state.setNPCInt(npcSpecialAction),
		"NPC_GetSpecialAction": state.getNPCInt(npcSpecialAction), "NPC_SetKeys": state.setNPCKeys,
		"NPC_GetKeys": state.getNPCKeys, "NPC_MeleeAttack": state.startNPCMelee,
		"NPC_StopMeleeAttack": state.stopNPCMelee, "NPC_IsMeleeAttacking": state.isNPCMelee,
		"NPC_EnableReloading": state.setNPCBool(npcReloadEnabled), "NPC_IsReloadEnabled": state.getNPCBool(npcReloadEnabled),
		"NPC_IsReloading": state.getNPCBool(npcReloading), "NPC_EnableInfiniteAmmo": state.setNPCBool(npcInfiniteAmmo),
		"NPC_IsInfiniteAmmoEnabled": state.getNPCBool(npcInfiniteAmmo), "NPC_Shoot": state.shootNPC,
		"NPC_IsShooting": state.getNPCBool(npcShooting), "NPC_AimAt": state.aimNPC,
		"NPC_AimAtPlayer": state.aimNPCAtPlayer, "NPC_StopAim": state.stopNPCAim,
		"NPC_IsAiming": state.getNPCBool(npcAiming), "NPC_IsAimingAtPlayer": state.isNPCAimingAtPlayer,
		"NPC_GetPlayerAimingAt": state.getNPCAimPlayer, "NPC_SetWeaponAccuracy": state.setNPCWeaponFloat(npcWeaponAccuracy),
		"NPC_GetWeaponAccuracy": state.getNPCWeaponFloat(npcWeaponAccuracy), "NPC_SetWeaponReloadTime": state.setNPCWeaponInt(npcWeaponReloadTime),
		"NPC_GetWeaponReloadTime": state.getNPCWeaponInt(npcWeaponReloadTime), "NPC_GetWeaponActualReloadTime": state.getNPCWeaponInt(npcWeaponReloadTime),
		"NPC_SetWeaponShootTime": state.setNPCWeaponInt(npcWeaponShootTime), "NPC_GetWeaponShootTime": state.getNPCWeaponInt(npcWeaponShootTime),
		"NPC_SetWeaponClipSize": state.setNPCWeaponInt(npcWeaponClipSize), "NPC_GetWeaponClipSize": state.getNPCWeaponInt(npcWeaponClipSize),
		"NPC_GetWeaponActualClipSize": state.getNPCWeaponInt(npcWeaponClipSize), "NPC_SetWeaponSkillLevel": state.setNPCWeaponInt(npcWeaponSkills),
		"NPC_GetWeaponSkillLevel": state.getNPCWeaponInt(npcWeaponSkills), "NPC_ResetAnimation": state.clearNPCAnimation,
		"NPC_ClearAnimations": state.clearNPCAnimation, "NPC_SetAnimation": state.setNPCAnimation,
		"NPC_GetAnimation": state.getNPCAnimation, "NPC_ApplyAnimation": state.applyNPCAnimation,
	}
}

func npcWeapon(npc *testNPC) *int         { return &npc.weapon }
func npcAmmo(npc *testNPC) *int           { return &npc.ammo }
func npcAmmoInClip(npc *testNPC) *int     { return &npc.ammoInClip }
func npcFightingStyle(npc *testNPC) *int  { return &npc.fightingStyle }
func npcWeaponState(npc *testNPC) *int    { return &npc.weaponState }
func npcSpecialAction(npc *testNPC) *int  { return &npc.specialAction }
func npcReloadEnabled(npc *testNPC) *bool { return &npc.reloadEnabled }
func npcReloading(npc *testNPC) *bool     { return &npc.reloading }
func npcInfiniteAmmo(npc *testNPC) *bool  { return &npc.infiniteAmmo }
func npcShooting(npc *testNPC) *bool      { return &npc.shooting }
func npcAiming(npc *testNPC) *bool        { return &npc.aiming }

func npcWeaponAccuracy(npc *testNPC) map[int]float32 { return npc.weaponAccuracy }
func npcWeaponReloadTime(npc *testNPC) map[int]int   { return npc.weaponReloadTime }
func npcWeaponShootTime(npc *testNPC) map[int]int    { return npc.weaponShootTime }
func npcWeaponClipSize(npc *testNPC) map[int]int     { return npc.weaponClipSize }
func npcWeaponSkills(npc *testNPC) map[int]int       { return npc.weaponSkills }

func (state *npcState) setNPCBool(field func(*testNPC) *bool) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		npc, ok := state.npc(params)
		if !ok {
			return 0, nil
		}

		*field(npc) = len(params) < 2 || params[1] != 0

		return 1, nil
	}
}

func (state *npcState) getNPCBool(field func(*testNPC) *bool) backend.NativeFunc {
	return func(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		npc, ok := state.npc(params)
		if ok && *field(npc) {
			return 1, nil
		}

		return 0, nil
	}
}

func (state *npcState) setNPCKeys(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok || len(params) < 4 {
		return 0, nil
	}

	npc.keys = [3]int{int(params[1]), int(params[2]), int(params[3])}

	return 1, nil
}

func (state *npcState) getNPCKeys(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok || len(params) < 4 {
		return 0, nil
	}

	return writeVehicleInts(ctx, params[1:4], npc.keys[:])
}

func (state *npcState) startNPCMelee(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok {
		return 0, nil
	}

	npc.melee = true

	return 1, nil
}

func (state *npcState) stopNPCMelee(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if !ok {
		return 0, nil
	}

	npc.melee = false

	return 1, nil
}

func (state *npcState) isNPCMelee(_ backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
	npc, ok := state.npc(params)
	if ok && npc.melee {
		return 1, nil
	}

	return 0, nil
}
