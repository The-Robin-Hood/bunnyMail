import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Field,
  FieldDescription,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";

import { Info, ChevronRight } from "lucide-react";
import { toast } from "sonner";
import { z } from "zod";

const AccountSetupStep1Schema = z.object({
  fullName: z.string().min(3, "Don't be shy, enter your full name"),
  email: z.email("Wait, that doesn't look like a valid email address"),
  password: z.string().min(6, "Seriously? Your email provider allows that short of a password?"),
  rememberPassword: z.string()
});

type AccountSetupStep1Data = z.infer<typeof AccountSetupStep1Schema>;

export default function AccountSetupStep1({
  formData,
  updateField,
  handleNext,
}: {
  formData: AccountSetupStep1Data;
  updateField: (field: string, value: string) => void;
  handleNext: () => void;
}) {
  const validateFields = () => {
    const result = AccountSetupStep1Schema.safeParse(formData);
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
        <Field className="gap-1">
          <FieldLabel>Full Name</FieldLabel>
          <Input
            type="text"
            placeholder="Robin Hood"
            value={formData.fullName}
            onChange={(e) => updateField("fullName", e.target.value)}
          />
        </Field>

        <Field className="gap-1">
          <FieldLabel>
            <span>Email Address</span>
            <Tooltip>
              <TooltipTrigger asChild>
                <Info className="size-4" />
              </TooltipTrigger>
              <TooltipContent>
                <p>Auto-configures for Gmail, Yahoo, Outlook</p>
              </TooltipContent>
            </Tooltip>
          </FieldLabel>
          <Input
            type="email"
            placeholder="robin.hood@example.com"
            value={formData.email}
            onChange={(e) => updateField("email", e.target.value)}
          />
          <FieldDescription className="text-xs">
            Server settings will be auto-configured for known providers.
          </FieldDescription>
        </Field>

        <Field className="gap-1">
          <FieldLabel>Password</FieldLabel>
          <Input
            type="password"
            placeholder="Password"
            value={formData.password}
            onChange={(e) => updateField("password", e.target.value)}
          />
          <FieldDescription className="mt-1 text-xs">
            Your password will be stored locally on your device.
          </FieldDescription>
        </Field>

        <Field className="flex-row items-center gap-2">
          <Checkbox
            className="max-w-4"
            id="remember-password"
            value={formData.rememberPassword ? "true" : "false"}
            onCheckedChange={(checked: boolean) => {
              updateField("rememberPassword", String(checked));
            }}
          />
          <FieldLabel htmlFor="remember-password" className="mb-0">
            Remember Password
          </FieldLabel>
        </Field>
      </FieldGroup>

      <div className="mt-6 flex gap-3">
        <Button
          onClick={() => validateFields() && handleNext()}
          className="w-full"
        >
          Next <ChevronRight className="ml-2 w-4 h-4" />
        </Button>
      </div>
    </div>
  );
}
