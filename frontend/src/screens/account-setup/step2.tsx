import { Button } from "@/components/ui/button";
import { Field, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { ChevronLeft, ChevronRight } from "lucide-react";
import { toast } from "sonner";
import { z } from "zod";

const AccountSetupStep2Schema = z.object({
  imapServer: z.string().min(1, "IMAP Server is required"),
  imapPort: z.string().min(1, "IMAP Port is required"),
  imapSecurity: z.string(),
  imapAuth: z.string(),
  smtpServer: z.string().min(1, "SMTP Server is required"),
  smtpPort: z.string().min(1, "SMTP Port is required"),
  smtpSecurity: z.string(),
  smtpAuth: z.string()
});

type AccountSetupStep2Data = z.infer<typeof AccountSetupStep2Schema>;

export default function AccountSetupStep2({
  formData,
  updateField,
  handleNext,
  handleBack,
}: {
  formData: AccountSetupStep2Data;
  updateField: (field: string, value: string) => void;
  handleNext: () => void;
  handleBack: () => void;
}) {
  const validateFields = () => {
    const result = AccountSetupStep2Schema.safeParse(formData);
    if (!result.success) {
      const firstError = result.error.issues[0];
      toast.error(firstError.message);
      return false;
    }
    return true;
  };

  return (
    <div>
      <FieldGroup className="gap-3">
        <div className="space-y-3">
          <Field className="gap-1">
            <FieldLabel>IMAP Server</FieldLabel>
            <Input
              type="text"
              placeholder="imap.gmail.com"
              value={formData.imapServer}
              onChange={(e) => updateField("imapServer", e.target.value)}
            />
          </Field>
          <div className="grid grid-cols-2 gap-3">
            <Field className="gap-1">
              <FieldLabel>Port</FieldLabel>
              <Input
                type="text"
                placeholder="993"
                value={formData.imapPort}
                onChange={(e) => updateField("imapPort", e.target.value)}
              />
            </Field>
            <Field className="gap-1">
              <FieldLabel>Connection Security</FieldLabel>
              <Select
                value={formData.imapSecurity}
                onValueChange={(val) => updateField("imapSecurity", val)}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Select security" />
                </SelectTrigger>
                <SelectContent>
                  <SelectGroup>
                    <SelectItem value="none">None</SelectItem>
                    <SelectItem value="start_tls">STARTTLS</SelectItem>
                    <SelectItem value="ssl_tls">SSL/TLS</SelectItem>
                  </SelectGroup>
                </SelectContent>
              </Select>
            </Field>
          </div>
          <Field className="gap-1">
            <FieldLabel>Authentication Method</FieldLabel>
            <Select
              value={formData.imapAuth}
              onValueChange={(val) => updateField("imapAuth", val)}
            >
              <SelectTrigger>
                <SelectValue placeholder="Select authentication" />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectItem value="Normal Password">
                    Normal Password
                  </SelectItem>
                  <SelectItem value="OAuth2">OAuth2</SelectItem>
                </SelectGroup>
              </SelectContent>
            </Select>
          </Field>
        </div>
        <div className="border-t my-1"></div>
        <div className="space-y-3">
          <Field className="gap-1">
            <FieldLabel>SMTP Server</FieldLabel>
            <Input
              type="text"
              placeholder="smtp.gmail.com"
              value={formData.smtpServer}
              onChange={(e) => updateField("smtpServer", e.target.value)}
            />
          </Field>
          <div className="grid grid-cols-2 gap-3">
            <Field className="gap-1">
              <FieldLabel>Port</FieldLabel>
              <Input
                type="text"
                placeholder="587"
                value={formData.smtpPort}
                onChange={(e) => updateField("smtpPort", e.target.value)}
              />
            </Field>
            <Field className="gap-1">
              <FieldLabel>Connection Security</FieldLabel>
              <Select
                value={formData.smtpSecurity}
                onValueChange={(val) => updateField("smtpSecurity", val)}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Select security" />
                </SelectTrigger>
                <SelectContent>
                  <SelectGroup>
                    <SelectItem value="none">None</SelectItem>
                    <SelectItem value="start_tls">STARTTLS</SelectItem>
                    <SelectItem value="ssl_tls">SSL/TLS</SelectItem>
                  </SelectGroup>
                </SelectContent>
              </Select>
            </Field>
          </div>
          <Field className="gap-1">
            <FieldLabel>Authentication Method</FieldLabel>
            <Select
              value={formData.smtpAuth}
              onValueChange={(val) => updateField("smtpAuth", val)}
            >
              <SelectTrigger>
                <SelectValue placeholder="Select authentication" />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectItem value="Normal Password">
                    Normal Password
                  </SelectItem>
                  <SelectItem value="OAuth2">OAuth2</SelectItem>
                </SelectGroup>
              </SelectContent>
            </Select>
          </Field>
        </div>
      </FieldGroup>

      <div className="mt-6 flex gap-3">
        <Button
          variant="outline"
          onClick={handleBack}
          className="flex-1 items-center justify-center"
        >
          <ChevronLeft className="w-5 h-5" /> <span>Back</span>
        </Button>
        <Button
          onClick={() => validateFields() && handleNext()}
          className="flex-1 items-center justify-center"
        >
          <span>Next</span>
          <ChevronRight className="w-5 h-5" />
        </Button>
      </div>
    </div>
  );
}
