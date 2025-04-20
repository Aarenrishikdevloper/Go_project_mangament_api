package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"rishik.com/constants"
	"rishik.com/enums"
)

// code for generating random invite code
func GenerateInviteCode() string {
	return randomString(8)
}

// code for generating Task code

// generating random code
func randomString(lenght int) string {
	const chrarset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, lenght)
	for i := range b {
		b[i] = chrarset[seededRand.Intn(len(chrarset))]
	}
	return string(b)
}

// var Unauthorize = errors.New("you do not have the necessary permissions to perform this action")
func RoleGuard(role enums.RoleType, requiredPermission []enums.PermissionType) error {
	permissiion, ok := constants.RolePermissions[role]
	if !ok {
		return errors.New("unauthorized: role does not exists")
	}
	permissionset := make(map[enums.PermissionType]struct{})
	for _, p := range permissiion {
		permissionset[p] = struct{}{}
	}
	for _, required := range requiredPermission {
		if _, exists := permissionset[required]; !exists {
			return errors.New("unauthorized:missing required permission")
		}
	}
	return nil
}
func GenerateTaskCode() string {
	id := uuid.New().String()
	return fmt.Sprintf("task-%s", strings.ReplaceAll(id, "-", "")[:3])
}
