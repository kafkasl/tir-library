package models

//import (
//	"sort"
//	"testing"
//)
//
//func sortAccounts(Accounts []Account) {
//	sort.Slice(Accounts, func(i, j int) bool {
//		return Accounts[i].ID < Accounts[j].ID
//	})
//}
//
//// Creates new Accounts into the DB
//func TestCreateAccounts(t *testing.T) {
//
//	prepareTestDatabase()
//
//	expectedAccounts := make([]Account, 1)
//	expectedAccounts[0] = Account{Email: "arthur@heart.gold", Password: "424242", Admin: false}
//	expectedAccounts[1] = Account{Email: "admin@superuser.org", Password: "lens.amelie", Admin: true}
//
//	for _, u := range expectedAccounts {
//		u.Create(u.Admin)
//	}
//	gotAccounts := GetUsers()
//
//	found := 0
//	for _, got := range gotAccounts {
//
//		for _, expected := range expectedAccounts {
//
//			if expected.Admin == got.Admin && expected.Email == got.Email && expected.Token == got.Token {
//				found++
//				break
//			}
//		}
//
//	}
//	if found != len(expectedAccounts) {
//		t.Errorf("Inserted registers %v not found with GetUsers.", gotAccounts)
//	}
//}
//
//func TestGetAccounts(t *testing.T) {
//
//	prepareTestDatabase()
//
//	expectedAccounts := make([]Account, 2)
//	expectedAccounts[0] = Account{Email: "julian.delphiki@battle.school", Password: "beanie"}
//	expectedAccounts[1] = Account{Email: "trurl@cyberiad.cosmos", Password: "Klapaucius"}
//
//	gotAccounts := GetUsers()
//
//	if len(gotAccounts) < len(expectedAccounts) {
//		t.Errorf("Expected at least %v records, got: %v\n%v\n", len(expectedAccounts), len(gotAccounts), gotAccounts)
//		return
//	}
//
//	sortAccounts(expectedAccounts)
//	sortAccounts(gotAccounts)
//
//	for i := 0; i < len(expectedAccounts); i++ {
//
//		got, expected := gotAccounts[i], expectedAccounts[i]
//
//		if expected.Admin != got.Admin || expected.Email != got.Email || expected.Token != got.Token {
//			t.Errorf("Handler returned wrong Account: got \n %v\n %v\n %v\n "+
//				"want\n %v\n %v\n %v\n",
//				got.Email, got.Token, got.Admin, expected.Email, expected.Token, expected.Admin)
//		}
//
//	}
//
//}
