using System;
using System.Security.Cryptography;

namespace Agent.Core.Utils
{
    public static class ECDHUtils
    {
        public static (byte[] privateKey, byte[] publicKey) GenerateECDHKeyPair()
        {
            using (var ecdh = ECDiffieHellman.Create(ECCurve.NamedCurves.nistP256))
            {
                var parameters = ecdh.ExportParameters(true);

                // Private key (D) - 32 bytes
                byte[] privateKey = parameters.D;

                // Create public key in correct format (0x04 || X || Y)
                byte[] publicKey = new byte[65];
                publicKey[0] = 0x04; // Uncompressed point format
                Array.Copy(parameters.Q.X, 0, publicKey, 1, 32);
                Array.Copy(parameters.Q.Y, 0, publicKey, 33, 32);

                return (privateKey, publicKey);
            }
        }

        public static byte[] DeriveECDHSharedSecret(byte[] privateKeyBytes, byte[] peerPublicKeyBytes)
        {
            using (var ecdh = ECDiffieHellman.Create(ECCurve.NamedCurves.nistP256))
            {
                // Import private key
                var parameters = new ECParameters
                {
                    Curve = ECCurve.NamedCurves.nistP256,
                    D = privateKeyBytes
                };
                ecdh.ImportParameters(parameters);

                // Import peer's public key
                using (var peerKey = ECDiffieHellman.Create(ECCurve.NamedCurves.nistP256))
                {
                    var peerParams = new ECParameters
                    {
                        Curve = ECCurve.NamedCurves.nistP256,
                        Q = new ECPoint
                        {
                            X = peerPublicKeyBytes.AsSpan(1, 32).ToArray(),
                            Y = peerPublicKeyBytes.AsSpan(33, 32).ToArray()
                        }
                    };
                    peerKey.ImportParameters(peerParams);

                    // Derive shared secret
                    byte[] sharedSecret = ecdh.DeriveKeyMaterial(peerKey.PublicKey);

                    // Hash the shared secret to get the final key
                    using (SHA256 sha256 = SHA256.Create())
                    {
                        return sha256.ComputeHash(sharedSecret);
                    }
                }
            }
        }
    }
}