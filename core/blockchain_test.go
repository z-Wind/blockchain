package core

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestBlockchain_NewBlockchain(t *testing.T) {
	bc := NewBlockchain()

	if len(bc.Chain) != 1 {
		t.Error(bc.Chain)
	}
}

func Test_block_proofWork(t *testing.T) {
	tests := []struct {
		name string
		b    *Block
	}{
		// TODO: Add test cases.
		{"easy1", &Block{Index: 1, Timestamp: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), Transactions: []*Transaction{}, Proof: 1, PreviousHash: "123"}},
		{"easy2", &Block{Index: 2, Timestamp: time.Date(2019, time.November, 10, 23, 0, 0, 0, time.UTC), Transactions: []*Transaction{}, Proof: 1, PreviousHash: "123365544"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.proofWork()
		})
	}
}

func TestBlockchain_NewTransaction(t *testing.T) {
	type args struct {
		sender    string
		recipient string
		amount    float64
	}
	bc := NewBlockchain()
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{"t1", args{"sender", "recipient", 123}},
		{"t2", args{"sender", "recipient", -123}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bc.NewTransaction(tt.args.sender, tt.args.recipient, tt.args.amount)
		})
	}
}

func TestBlockchain_RegisterNode(t *testing.T) {
	type args struct {
		address string
	}
	bc := NewBlockchain()
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"a1", args{"http:11.11.111.11"}, "http:11.11.111.11"},
		{"a2", args{"http:11.11.111.11"}, "http:11.11.111.11"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bc.RegisterNode(tt.args.address)
			got := tt.args.address
			if _, ok := bc.Nodes[tt.want]; !ok {
				t.Errorf("block.registerNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlockchain_validChain(t *testing.T) {
	// PASS
	bc := NewBlockchain()
	bc.NewTransaction("sender", "recipient", 123)
	bc.NewTransaction("sender", "recipient", 456)
	bc.NewBlock()

	if ok, err := bc.validChain(); !ok || err != nil {
		t.Errorf("Blockchain.validChain() should be PASS\nbc = %+v, error: %s", bc, err)
	}

	// NG
	bc = NewBlockchain()
	bc.NewTransaction("sender", "recipient", 123)
	bc.NewTransaction("sender", "recipient", 456)
	bc.NewBlock()
	bc.Chain[1].Proof = 0

	if ok, err := bc.validChain(); ok {
		t.Errorf("Blockchain.validChain() should be failed\n bc = %+v, error: %s", bc, err)
	}

	// NG
	bc = NewBlockchain()
	bc.NewTransaction("sender", "recipient", 123)
	bc.NewTransaction("sender", "recipient", 456)
	bc.NewBlock()
	bc.Chain[1].PreviousHash = "123"
	bc.Chain[1].proofWork()

	if ok, err := bc.validChain(); ok {
		t.Errorf("Blockchain.validChain() should be failed\n bc = %+v, error: %s", bc, err)
	}

}

func TestBlockchain_resolveConflicts(t *testing.T) {
	bc1 := NewBlockchain()
	bc1.NewTransaction("sender", "recipient", 123)
	bc1.NewTransaction("sender", "recipient", 456)
	bc1.NewBlock()

	bc2 := NewBlockchain()
	bc2.NewTransaction("sender", "recipient", 123)
	bc2.NewTransaction("sender", "recipient", 456)
	bc2.NewBlock()
	bc2.NewTransaction("sender", "recipient", 123)
	bc2.NewTransaction("sender", "recipient", 456)
	bc2.NewBlock()

	changed, err := bc1.ResolveConflicts(bc2)
	if err != nil {
		t.Errorf("error = %s", err)
	}
	if want := true; changed != want {
		t.Errorf("block.resolveConflicts() = %v, want %v", changed, want)
	}

	changed, err = bc1.ResolveConflicts(bc2)
	if err != nil {
		t.Errorf("error = %s", err)
	}
	if want := false; changed != want {
		t.Errorf("block.resolveConflicts() = %v, want %v", changed, want)
	}
}

func Test_block_String(t *testing.T) {
	tests := []struct {
		name string
		b    *Block
		want string
	}{
		// TODO: Add test cases.
		{"easy1", &Block{Index: 1, Timestamp: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), Transactions: []*Transaction{}, Proof: 1, PreviousHash: "123"},
			"{\n#1 2009-11-10 23:00:00 +0000 UTC\ntransactions:[]\nproof:1\nprevious Hash:123\n}"},
		{"easy2", &Block{Index: 1, Timestamp: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), Transactions: []*Transaction{{"sender", "recipient", 123}, {"sender", "recipient", 123}}, Proof: 1, PreviousHash: "123"},
			"{\n#1 2009-11-10 23:00:00 +0000 UTC\ntransactions:[{s:sender, r:recipient, $123.00} {s:sender, r:recipient, $123.00}]\nproof:1\nprevious Hash:123\n}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.String(); got != tt.want {
				t.Errorf("block.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlockchain_String(t *testing.T) {
	bc := NewBlockchain()
	bc.Chain[0].Timestamp = time.Date(2019, time.September, 7, 25, 9, 53, 0, time.UTC)
	bc.Chain[0].proofWork()
	bc.NewTransaction("sender", "recipient", 123)
	bc.NewTransaction("sender", "recipient", 456)
	bc.NewBlock()
	bc.Chain[1].Timestamp = time.Date(2019, time.September, 7, 25, 9, 53, 1, time.UTC)
	bc.Chain[1].proofWork()
	fmt.Print("==========================")

	want := "transactions:[{s:sender, r:recipient, $123.00} {s:sender, r:recipient, $456.00}]\nproof:"
	if got := bc.String(); !strings.Contains(got, want) {
		t.Errorf("Blockchain.String() = %s, want %s", got, want)
	}
}

func TestBlock_validProof(t *testing.T) {
	tests := []struct {
		name    string
		b       *Block
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
		{"easy", &Block{Index: 1, Timestamp: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), Transactions: []*Transaction{}, Proof: 1, PreviousHash: "123"}, false, false},
		{"check", &Block{Index: 11, Timestamp: time.Date(2019, time.November, 10, 23, 0, 0, 0, time.UTC), Transactions: []*Transaction{}, Proof: 1, PreviousHash: "123"}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.validProof()
			if (err != nil) != tt.wantErr {
				t.Errorf("Block.validProof() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Block.validProof() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlock_hash(t *testing.T) {
	tests := []struct {
		name    string
		b       *Block
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"easy1", &Block{Index: 1, Timestamp: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), Transactions: []*Transaction{}, Proof: 1, PreviousHash: "123"}, "e417bc7b35fbc7c2193e8f743604337e4bc1337150da1ce90f43a48551891ead", false},
		{"easy2", &Block{Index: 2, Timestamp: time.Date(2019, time.November, 10, 23, 0, 0, 0, time.UTC), Transactions: []*Transaction{}, Proof: 1, PreviousHash: "123365544"}, "a190dfe17f5454e4ff884972071579d7d9f975bea41daf224baf1323e4c8276b", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.hash()
			if (err != nil) != tt.wantErr {
				t.Errorf("Block.hash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Block.hash() = %v, want %v", got, tt.want)
			}
		})
	}
}
