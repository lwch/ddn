package dht

import "crypto/rand"

const ver = "1000"

func randID() [20]byte {
	const charMap = "0123456789abcdef"
	var id [20]byte
	id[0] = '-'
	id[1] = 'D'
	id[2] = 'N'
	copy(id[3:], ver)
	id[7] = '-'
	for i := 8; i < 20; {
		n, err := rand.Read(id[i:])
		if err != nil {
			for j := i; j < 20; j++ {
				id[j] = 'f'
			}
			return id
		}
		for j := i; j < i+n; j++ {
			id[j] = charMap[int(id[j])%len(charMap)]
		}
		i += n
	}
	return id
}
