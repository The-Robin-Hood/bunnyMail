"use client";

import * as React from "react";
import { ArchiveX, Inbox, Send, Trash2 } from "lucide-react";

import { NavUser } from "@/components/nav-user";
import { Label } from "@/components/ui/label";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarHeader,
  SidebarInput,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { Switch } from "@/components/ui/switch";
import { G_GetMessagesByAccount } from "@/../wailsjs/go/main/App";
import type { model } from "@/../wailsjs/go/models";
import { useMailStore } from "@/store/mailStore";
import { cn, convertTime, truncateText } from "@/lib/utils";

const folders = [
  {
    title: "Inbox",
    url: "#",
    icon: Inbox,
    isActive: true,
  },
  {
    title: "Sent",
    url: "#",
    icon: Send,
    isActive: false,
  },
  {
    title: "Junk",
    url: "#",
    icon: ArchiveX,
    isActive: false,
  },
  {
    title: "Trash",
    url: "#",
    icon: Trash2,
    isActive: false,
  },
];

export function AppSidebar() {
  const [activeItem, setActiveItem] = React.useState(folders[0]);
  const [mails, setMails] = React.useState<model.Message[]>([]);
  const setSelectedMail = useMailStore((state) => state.selectMessage);
  const selectedMail = useMailStore((state) => state.selectedMessage);
  const selectedAccount = useMailStore((state) => state.selectedAccount);

  React.useEffect(() => {
    if (!selectedAccount) return;
    G_GetMessagesByAccount(selectedAccount.id, 20)
      .then((messages) => {
        console.log("Fetched messages in sidebar:", messages);
        setMails(messages);
      })
      .catch((err) => {
        console.error("Fetching messages failed in sidebar:", err);
      });
  }, [selectedAccount]);

  return (
    <>
      <Sidebar
        collapsible="none"
        className="w-[calc(var(--sidebar-width-icon)+1px)]! border-r h-screen"
      >
        <SidebarHeader>
          <SidebarMenu>
            <SidebarMenuItem>
              <img
                src="/logo 1.png"
                alt="Bunny Mail"
                className="h-8 w-12 rounded"
              />
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarHeader>
        <SidebarContent>
          <SidebarGroup>
            <SidebarGroupContent className="px-0">
              <SidebarMenu>
                {folders.map((item) => (
                  <SidebarMenuItem key={item.title}>
                    <SidebarMenuButton
                      tooltip={{
                        children: item.title,
                        hidden: false,
                      }}
                      onClick={() => {
                        setActiveItem(item);
                      }}
                      isActive={activeItem?.title === item.title}
                      className="px-2"
                    >
                      <item.icon />
                      <span>{item.title}</span>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>
        <SidebarFooter>
          <NavUser />
        </SidebarFooter>
      </Sidebar>

      <Sidebar
        collapsible="none"
        className="flex h-screen border-r min-w-75 max-w-sm"
      >
        <SidebarHeader className="gap-3.5 border-b p-4">
          <div className="flex w-full items-center justify-between">
            <div className="text-foreground text-base font-medium">
              {activeItem?.title}
            </div>
            <Label className="flex items-center gap-2 text-sm">
              <span>Unreads</span>
              <Switch className="shadow-none" />
            </Label>
          </div>
          <SidebarInput placeholder="Type to search..." />
        </SidebarHeader>
        <SidebarContent>
          <SidebarGroup className="px-0">
            <SidebarGroupContent>
              {mails.map((mail) => (
                <div
                  key={mail.id}
                  className={cn(
                    "hover:bg-sidebar-accent hover:text-sidebar-accent-foreground flex flex-col items-start gap-2 border-b p-4 text-sm leading-tight whitespace-nowrap last:border-b-0",
                    selectedMail && selectedMail.id === mail.id
                      ? "bg-sidebar-accent text-sidebar-accent-foreground"
                      : "text-foreground",
                  )}
                  onClick={(e) => {
                    e.preventDefault();
                    setSelectedMail(mail);
                  }}
                >
                  <div className="flex w-full items-center gap-2">
                    <span className="font-semibold">{mail.from_name}</span>{" "}
                    <span className="ml-auto text-xs">
                      {convertTime(mail.received_at)}
                    </span>
                  </div>
                  <span className="font-normal text-xs text-html">
                    {truncateText(mail.subject || "(No Subject)", 45)}
                  </span>
                </div>
              ))}
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>
      </Sidebar>
    </>
  );
}
