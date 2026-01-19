import { AppSidebar } from "@/components/app-sidebar";
import {
  SidebarInset,
  SidebarProvider,
} from "@/components/ui/sidebar";
import { useMailStore } from "@/store/mailStore";

export default function Home() {
  const selectedMail = useMailStore((state) => state.selectedMessage);

  return (
    <SidebarProvider style={{"--sidebar-width": "350px"} as React.CSSProperties}>
      <AppSidebar/>
      <SidebarInset>
        <div className="flex flex-1 flex-col gap-4 p-4 bg-primary-foreground ">
          {selectedMail ? (
            <div className="flex flex-1 flex-col h-full">
              <h1 className="text-2xl font-bold mb-4">
                {selectedMail.subject || "(No Subject)"}
              </h1>
              <p>
                <strong>From:</strong> {selectedMail.from_name} &lt;
                {selectedMail.from_address}&gt;
              </p>
              <p>
                <strong>To:</strong> {selectedMail.to_addresses}
              </p>
              <p>
                <strong>Received At:</strong>{" "}
                {new Date(selectedMail.received_at).toLocaleString()}
              </p>
              <hr className="my-4" />
              <iframe
                srcDoc={selectedMail.body_html}
                sandbox="allow-same-origin allow-popups"
                className="h-full w-full border bg-html"
                title="Email content"
              />
            </div>
          ) : (
            <div className="flex flex-1 items-center justify-center">
              <p className="text-muted-foreground">
                Select an email to view its details.
              </p>
            </div>
          )}
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
}
