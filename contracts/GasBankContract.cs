using Neo;
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Native;
using Neo.SmartContract.Framework.Services;
using System;
using System.ComponentModel;
using System.Numerics;

namespace NeoOracle
{
    [DisplayName("GasBank")]
    [ManifestExtra("Author", "Neo Oracle Team")]
    [ManifestExtra("Email", "contact@neo-oracle.io")]
    [ManifestExtra("Description", "Neo N3 Gas Bank Contract")]
    public class GasBankContract : SmartContract
    {
        // Events
        [DisplayName("Deposit")]
        public static event Action<UInt160, BigInteger> OnDeposit;
        
        [DisplayName("Withdraw")]
        public static event Action<UInt160, BigInteger> OnWithdraw;
        
        [DisplayName("Allocation")]
        public static event Action<UInt160, string, BigInteger> OnAllocation;
        
        // Storage keys
        private static readonly byte[] BalancePrefix = new byte[] { 0x01 };
        private static readonly byte[] AuthorizedOracleKey = new byte[] { 0x02 };
        
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
        
        // Deposit GAS into the bank
        public static bool Deposit()
        {
            // Get the transaction sender
            UInt160 sender = Runtime.CallingScriptHash;
            
            // Get the GAS payment amount
            BigInteger amount = GetGasPayment();
            if (amount <= 0)
                throw new Exception("No GAS payment found");
                
            // Update user balance
            BigInteger currentBalance = GetBalance(sender);
            BigInteger newBalance = currentBalance + amount;
            SetBalance(sender, newBalance);
            
            // Emit deposit event
            OnDeposit(sender, amount);
            
            return true;
        }
        
        // Withdraw GAS from the bank
        public static bool Withdraw(BigInteger amount)
        {
            // Validate inputs
            if (amount <= 0)
                throw new Exception("Invalid amount");
                
            // Get the sender
            UInt160 sender = Runtime.CheckWitness();
            if (sender == null)
                throw new Exception("Not authorized");
                
            // Check balance
            BigInteger currentBalance = GetBalance(sender);
            if (currentBalance < amount)
                throw new Exception("Insufficient balance");
                
            // Update balance
            BigInteger newBalance = currentBalance - amount;
            SetBalance(sender, newBalance);
            
            // Transfer GAS to the sender
            GAS.Transfer(Runtime.ExecutingScriptHash, sender, amount);
            
            // Emit withdrawal event
            OnWithdraw(sender, amount);
            
            return true;
        }
        
        // Allocate GAS for a function execution
        public static bool AllocateGas(UInt160 user, string functionId, BigInteger amount)
        {
            // Only authorized oracle can allocate gas
            UInt160 callingScript = Runtime.CallingScriptHash;
            if (callingScript != GetAuthorizedOracle())
                throw new Exception("Not authorized");
                
            // Validate inputs
            if (amount <= 0)
                throw new Exception("Invalid amount");
                
            // Check user balance
            BigInteger currentBalance = GetBalance(user);
            if (currentBalance < amount)
                throw new Exception("Insufficient balance");
                
            // Update balance
            BigInteger newBalance = currentBalance - amount;
            SetBalance(user, newBalance);
            
            // Emit allocation event
            OnAllocation(user, functionId, amount);
            
            return true;
        }
        
        // Get user's balance
        public static BigInteger GetBalance(UInt160 user)
        {
            StorageMap balanceMap = new(Storage.CurrentContext, BalancePrefix);
            byte[] balance = balanceMap.Get(user);
            
            if (balance == null || balance.Length == 0)
                return 0;
                
            return (BigInteger)StdLib.Deserialize(balance);
        }
        
        // Set user's balance
        private static void SetBalance(UInt160 user, BigInteger amount)
        {
            StorageMap balanceMap = new(Storage.CurrentContext, BalancePrefix);
            balanceMap.Put(user, StdLib.Serialize(amount));
        }
        
        // Get owner address
        private static UInt160 GetOwner() => (UInt160)Storage.Get(Storage.CurrentContext, "owner");
        
        // Get authorized oracle address
        public static UInt160 GetAuthorizedOracle() => (UInt160)Storage.Get(AuthorizedOracleKey);
        
        // Helper to get GAS payment from current transaction
        private static BigInteger GetGasPayment()
        {
            Transaction tx = (Transaction)Runtime.ScriptContainer;
            foreach (TransactionOutput output in tx.GetOutputs())
            {
                if (output.AssetId == GAS.Hash && output.ScriptHash == Runtime.ExecutingScriptHash)
                {
                    return output.Value;
                }
            }
            return 0;
        }
    }
} 