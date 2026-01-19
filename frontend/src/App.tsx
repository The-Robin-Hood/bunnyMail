import { Toaster } from "sonner";
import { ThemeProvider } from "./components/theme-provider";
import AccountSetup from "./screens/account-setup";
import { useEffect, useState } from "react";
import Home from "./screens/home";
import { RadixLoader } from "./components/ui/loader";
import { useMailStore } from "./store/mailStore";

export default function App() {
  // const syncAccount = useMailStore((state) => state.syncAccount);
  const storedAccounts = useMailStore((state) => state.accounts);
  const loadAccounts = useMailStore((state) => state.loadAccounts);
  const [loading, setLoading] = useState<boolean>(true);

  async function sleep(ms: number) {
    return new Promise((resolve) => setTimeout(resolve, ms));
  }

  useEffect(() => {
    const load = async () => {
      await sleep(500);
      await loadAccounts();
      setLoading(false);
    };
    load();
  }, [loadAccounts]);

  return (
    <ThemeProvider defaultTheme="dark">
      {loading ? (
        <RadixLoader />
      ) : storedAccounts.length > 0 ? (
        <Home />
      ) : (
        <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
          <AccountSetup />
        </div>
      )}
      <Toaster position="bottom-center" />
    </ThemeProvider>
  );
}
