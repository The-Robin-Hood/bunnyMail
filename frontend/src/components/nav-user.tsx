import {
  BadgeCheck,
  Bell,
  ChevronsUpDown,
  CreditCard,
  LogOut,
  PlusCircle,
} from "lucide-react";

import { Avatar, AvatarImage } from "@/components/ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { createAvatar } from "@dicebear/core";
import { adventurerNeutral } from "@dicebear/collection";
import { useMailStore } from "@/store/mailStore";
import { cn } from "@/lib/utils";

const getAvatar = (seed: string) => {
  const avatar = createAvatar(adventurerNeutral, {
    seed: seed,
    backgroundColor: ["829ab1"],
    backgroundType: ["solid"],
    scale: 100,
  });
  return avatar.toDataUri();
};

export function NavUser() {
  const accounts = useMailStore((state) => state.accounts);
  const currentAccount = useMailStore((state) => state.selectedAccount);
  // const switchAccount = useMailStore((state) => state.selectAccount);

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              size="lg"
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground md:h-8 p-0"
            >
              <Avatar className="h-8 w-8 rounded-md justify-center items-center bg-transparent">
                <AvatarImage
                  className="h-8 w-8"
                  src={getAvatar(currentAccount?.name || "user")}
                  alt={currentAccount?.name}
                />
              </Avatar>
              <div className="grid flex-1 text-left text-sm leading-tight">
                <span className="truncate font-medium">
                  {currentAccount?.name}
                </span>
                <span className="truncate text-xs">
                  {currentAccount?.email}
                </span>
              </div>
              <ChevronsUpDown className="ml-auto size-4" />
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            className="w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg"
            side="right"
            align="end"
            sideOffset={4}
          >
            {/* Account Switcher - Clickable Label */}
            <DropdownMenuSub>
              <DropdownMenuSubTrigger className="p-0 font-normal">
                <div className="flex items-center gap-2 px-1 py-1.5 text-left text-sm w-full">
                  <Avatar className="h-8 w-8 rounded-lg">
                    <AvatarImage
                      src={getAvatar(currentAccount?.name || "user")}
                      alt={currentAccount?.name}
                    />
                  </Avatar>
                  <div className="grid flex-1 text-left text-sm leading-tight">
                    <span className="truncate font-medium">
                      {currentAccount?.name}
                    </span>
                    <span className="truncate text-xs">
                      {currentAccount?.email}
                    </span>
                  </div>
                </div>
              </DropdownMenuSubTrigger>
              <DropdownMenuSubContent className="min-w-56">
                {accounts?.map((account) => (
                  <DropdownMenuItem
                    className={cn(
                      account.id === currentAccount?.id
                        ? "bg-sidebar-accent text-sidebar-accent-foreground"
                        : "",
                    )}
                    key={account.email}
                    onClick={() => {
                      if (account.id === currentAccount?.id) return;
                      console.log(account);
                    }}
                  >
                    <div className="grid flex-1 text-left text-sm leading-tight">
                      <span className="truncate font-medium">
                        {account.name}
                      </span>
                      <span className="truncate text-xs text-muted-foreground">
                        {account.email}
                      </span>
                    </div>
                  </DropdownMenuItem>
                ))}
                <DropdownMenuSeparator />
                <DropdownMenuItem>
                  <PlusCircle className="mr-2 h-4 w-4" />
                  Add account
                </DropdownMenuItem>
              </DropdownMenuSubContent>
            </DropdownMenuSub>

            <DropdownMenuSeparator />

            <DropdownMenuGroup>
              <DropdownMenuItem>
                <BadgeCheck className="mr-2 h-4 w-4" />
                Account
              </DropdownMenuItem>
              <DropdownMenuItem>
                <CreditCard className="mr-2 h-4 w-4" />
                Billing
              </DropdownMenuItem>
              <DropdownMenuItem>
                <Bell className="mr-2 h-4 w-4" />
                Notifications
              </DropdownMenuItem>
            </DropdownMenuGroup>

            <DropdownMenuSeparator />

            <DropdownMenuItem>
              <LogOut className="mr-2 h-4 w-4" />
              Log out
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  );
}
