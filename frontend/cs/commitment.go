package cs

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/consensys/gnark/constraint/solver"
	"github.com/consensys/gnark/debug"
	"github.com/consensys/gnark/logger"
	"hash/fnv"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
)

func Bsb22CommitmentComputePlaceholder(mod *big.Int, input []*big.Int, output []*big.Int) error {
	if (len(os.Args) > 0 && (strings.HasSuffix(os.Args[0], ".test") || strings.HasSuffix(os.Args[0], ".test.exe"))) || debug.Debug {
		// usually we only run solver without prover during testing
		log := logger.Logger()
		log.Error().Msg("Augmented commitment hint not replaced. Proof will not be sound and verification will fail!")
		byteLen := 1 + mod.BitLen()/8
		toHash := make([]byte, byteLen*len(input)+32)
		for i, in := range input {
			inBytes := in.Bytes()
			copy(toHash[(i+1)*byteLen-len(inBytes):], inBytes)
		}
		if n, err := rand.Read(toHash[len(toHash)-32:]); err != nil || n != 32 {
			return fmt.Errorf("generated %d random bytes: %v", n, err)
		}
		hsh := sha256.New().Sum(toHash)
		output[0].SetBytes(hsh)
		output[0].Mod(output[0], mod)

		return nil
	}
	return fmt.Errorf("placeholder function: to be replaced by commitment computation")
}

var maxNbCommitments int
var m sync.Mutex

func RegisterBsb22CommitmentComputePlaceholder(index int) (hintId solver.HintID, err error) {

	hintName := "bsb22 commitment #" + strconv.Itoa(index)
	// mimic solver.GetHintID
	hf := fnv.New32a()
	if _, err = hf.Write([]byte(hintName)); err != nil {
		return
	}
	hintId = solver.HintID(hf.Sum32())

	m.Lock()
	if maxNbCommitments == index {
		maxNbCommitments++
		solver.RegisterNamedHint(Bsb22CommitmentComputePlaceholder, hintId)
	}
	m.Unlock()

	return
}
func init() {
	solver.RegisterHint(Bsb22CommitmentComputePlaceholder)
}
