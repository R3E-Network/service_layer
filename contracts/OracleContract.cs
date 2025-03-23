using Neo;
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Native;
using Neo.SmartContract.Framework.Services;
using System;
using System.ComponentModel;
using System.Numerics;

namespace NeoOracle
{
    [DisplayName("NeoOracle")]
    [ManifestExtra("Author", "Neo Oracle Team")]
    [ManifestExtra("Email", "contact@neo-oracle.io")]
    [ManifestExtra("Description", "Neo N3 Oracle Service Contract")]
    public class OracleContract : SmartContract
    {
        // Events
        [DisplayName("PriceUpdated")]
        public static event Action<string, BigInteger, BigInteger> OnPriceUpdated;
        
        [DisplayName("FunctionExecuted")]
        public static event Action<string, string> OnFunctionExecuted;
        
        // Storage keys
        private static readonly byte[] PricePrefix = new byte[] { 0x01 };
        private static readonly byte[] AuthorizedOracleKey = new byte[] { 0x02 };
        
        // Price data
        public class PriceData
        {
            public BigInteger Price;
            public BigInteger Timestamp;
        }
        
        // Contract methods
        
        // Initialize the contract
        public static void Initialize(UInt160 oracleAddress)
        {
            // Only contract owner can initialize
            if (!Runtime.CheckWitness(GetOwner()))
                throw new Exception("No authorization");
                
            // Check if already initialized
            if (GetAuthorizedOracle() != null)
                throw new Exception("Already initialized");
                
            // Store authorized oracle address
            Storage.Put(AuthorizedOracleKey, oracleAddress);
        }
        
        // Update price data
        public static void UpdatePrice(string token, BigInteger price, BigInteger timestamp)
        {
            // Only authorized oracle can update prices
            UInt160 callingScript = Runtime.CallingScriptHash;
            if (callingScript != GetAuthorizedOracle())
                throw new Exception("Not authorized");
                
            // Validate inputs
            if (price <= 0)
                throw new Exception("Invalid price");
                
            if (timestamp <= 0)
                throw new Exception("Invalid timestamp");
                
            // Create price data
            PriceData data = new PriceData
            {
                Price = price,
                Timestamp = timestamp
            };
            
            // Store in contract storage
            byte[] key = GetPriceKey(token);
            StorageMap priceMap = new(Storage.CurrentContext, PricePrefix);
            priceMap.Put(key, StdLib.Serialize(data));
            
            // Emit event
            OnPriceUpdated(token, price, timestamp);
        }
        
        // Get price data
        public static PriceData GetPrice(string token)
        {
            byte[] key = GetPriceKey(token);
            StorageMap priceMap = new(Storage.CurrentContext, PricePrefix);
            byte[] value = priceMap.Get(key);
            
            if (value == null || value.Length == 0)
                return null;
                
            return (PriceData)StdLib.Deserialize(value);
        }
        
        // Record function execution
        public static void RecordFunctionExecution(string functionId, string result)
        {
            // Only authorized oracle can record executions
            UInt160 callingScript = Runtime.CallingScriptHash;
            if (callingScript != GetAuthorizedOracle())
                throw new Exception("Not authorized");
                
            // Emit event
            OnFunctionExecuted(functionId, result);
        }
        
        // Get owner address
        private static UInt160 GetOwner() => (UInt160)Storage.Get(Storage.CurrentContext, "owner");
        
        // Get authorized oracle address
        public static UInt160 GetAuthorizedOracle() => (UInt160)Storage.Get(AuthorizedOracleKey);
        
        // Helper to construct price storage key
        private static byte[] GetPriceKey(string token) => token.ToByteArray();
    }
} 