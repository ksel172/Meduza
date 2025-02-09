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
                var parameters = new ECParameters
                {
                    Curve = ECCurve.NamedCurves.nistP256,
                    D = privateKeyBytes,
                    Q = new ECPoint
                    {
                        X = new byte[32],
                        Y = new byte[32]
                    }
                };
                ecdh.ImportParameters(parameters);

                var peerKey = ECDiffieHellman.Create();
                var peerParams = new ECParameters
                {
                    Curve = ECCurve.NamedCurves.nistP256,
                    Q = new ECPoint
                    {
                        X = peerPublicKeyBytes[1..33],
                        Y = peerPublicKeyBytes[33..65]
                    }
                };
                peerKey.ImportParameters(peerParams);

                return ecdh.DeriveKeyFromHash(peerKey.PublicKey, HashAlgorithmName.SHA256);
            }
        }
    }
}