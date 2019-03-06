package utils

func HashPublicKey(publicKey []byte) []byte {
	out := Base58Decode(publicKey)
	out = out[1 : len(out)-4]
	return out
}
