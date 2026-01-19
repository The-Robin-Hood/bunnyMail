import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableRow,
} from "@/components/ui/table";
import { Check, ChevronLeft } from "lucide-react";

export default function AccountSetupStep3({
  formData,
  handleBack,
  handleSubmit,
}: {
  formData: {
    fullName: string;
    email: string;
    password: string;
    rememberPassword: string;
    imapServer: string;
    imapPort: string;
    imapSecurity: string;
    imapAuth: string;
    smtpServer: string;
    smtpPort: string;
    smtpSecurity: string;
    smtpAuth: string;
  };
  handleBack: () => void;
  handleSubmit: () => void;
}) {
  console.log("Reviewing form data:", formData);
  return (
    <div>
      <div className="space-y-4">
        <div className="bg-accent p-4 rounded-lg space-y-3">
          <h3 className="font-semibold text-accent-foreground">
            Account Information
          </h3>
          <div className="grid grid-cols-2 gap-2 text-sm overflow-hidden">
            <span className="text-accent-foreground">Full Name:</span>
            <span className="font-medium">{formData.fullName || "—"}</span>
            <span className="text-accent-foreground">Email:</span>
            <span className="font-medium">{formData.email || "—"}</span>
            <span className="text-accent-foreground">Remember Password:</span>
            <span className="font-medium">
              {formData.rememberPassword ? "Yes" : "No"}
            </span>
          </div>
        </div>

        <div className="bg-accent p-4 rounded-lg space-y-3">
          <h3 className="font-semibold text-accent-foreground">
            Server Configuration
          </h3>
          <div className="space-y-3 text-sm">
            <Table>
              <TableBody>
                <TableRow>
                  <TableCell>Server</TableCell>
                  <TableCell>{formData.imapServer || "Not set"}</TableCell>
                  <TableCell>{formData.smtpServer || "Not set"}</TableCell>
                </TableRow>
                <TableRow>
                  <TableCell>Port</TableCell>
                  <TableCell>{formData.imapPort}</TableCell>
                  <TableCell>{formData.smtpPort}</TableCell>
                </TableRow>
                <TableRow>
                  <TableCell>Security</TableCell>
                  <TableCell>{formData.imapSecurity}</TableCell>
                  <TableCell>{formData.smtpSecurity}</TableCell>
                </TableRow>
                <TableRow>
                  <TableCell>Auth</TableCell>
                  <TableCell>{formData.imapAuth}</TableCell>
                  <TableCell>{formData.smtpAuth}</TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </div>
        </div>
      </div>

      <div className="mt-6 flex gap-3">
        <Button
          variant="outline"
          onClick={handleBack}
          className="flex-1 text-white"
        >
          <ChevronLeft className="mr-2 w-4 h-4" />
          <span>Back</span>
        </Button>
        <Button
          onClick={handleSubmit}
          className="flex-1 bg-green-600 hover:bg-green-700 text-white"
        >
          <Check className="mr-2 w-4 h-4" />
          <span>Complete Setup</span>
        </Button>
      </div>
    </div>
  );
}
