// src/store/mailStore.ts
import { create } from "zustand";
import { model } from "@wailsjs/go/models";
import {
  G_GetAccounts,
  G_GetMessagesByAccount,
  G_SyncAccount,
} from "@wailsjs/go/main/App";

interface MailState {
  // Accounts
  accounts: model.Account[];
  selectedAccount: model.Account | null;
  loadAccounts: () => Promise<void>;
  selectAccount: (account: model.Account) => void;

  // Messages
  messages: model.Message[];
  selectedMessage: model.Message | null;
  loadMessages: (accountId: number, limit?: number) => Promise<void>;
  selectMessage: (message: model.Message | null) => Promise<void>;

  // Sync
  isSyncing: boolean;
  syncAccount: (accountId: number) => Promise<void>;

  // Mailbox
  currentMailbox: string;
  setCurrentMailbox: (mailbox: string) => void;
}

export const useMailStore = create<MailState>()((set, get) => ({
  // Initial state
  accounts: [],
  selectedAccount: null,
  messages: [],
  selectedMessage: null,
  currentMailbox: "INBOX",
  isSyncing: false,

  selectAccount: (account) => {
    set({ selectedAccount: account, messages: [], selectedMessage: null });
    // Load messages for the selected account
    get().loadMessages(account.id, 50);
  },

  // Load accounts
  loadAccounts: async () => {
    const accounts = await G_GetAccounts();
    set({ accounts });
    if (accounts.length > 0) {
      get().selectAccount(accounts[0]);
    }
  },

  // Load messages
  loadMessages: async (accountId, limit = 50) => {
    const messages = await G_GetMessagesByAccount(accountId, limit);
    set({ messages });
  },

  // Select and mark as read
  selectMessage: async (message) => {
    set({ selectedMessage: message });
    if (message && !message.is_read) {
      // await G_MarkAsRead(message.id)
      // Update in list
      set((state) => ({
        messages: state.messages.map((m) =>
          m.id === message.id ? { ...m, is_read: true } : m,
        ),
      }));
    }
  },

  // Sync account
  syncAccount: async (accountId) => {
    set({ isSyncing: true });
    try {
      await G_SyncAccount(accountId, 20);
      await get().loadMessages(accountId, 20);
    } catch (error) {
      console.error("Sync failed:", error);
    } finally {
      set({ isSyncing: false });
    }
  },

  // Set current mailbox
  setCurrentMailbox: (mailbox) => {
    set({ currentMailbox: mailbox });
  },
}));
