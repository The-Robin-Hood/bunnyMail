import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useState } from "react";
import emailProviders from "@/const/EmailProviders";
import StepIndicator from "@/components/step-indicator";

import AccountSetupStep1 from "./step1";
import AccountSetupStep2 from "./step2";
import AccountSetupStep3 from "./step3";

import {G_AddAccount} from "@/../wailsjs/go/main/App" 

const initialFormData = {
  fullName: "",
  email: "",
  password: "",
  rememberPassword: "false",
  imapServer: "",
  imapPort: "993",
  imapSecurity: "ssl_tls",
  imapAuth: "Normal Password",
  smtpServer: "",
  smtpPort: "587",
  smtpSecurity: "ssl_tls",
  smtpAuth: "Normal Password",
};

export default function AccountSetup() {
  const [currentStep, setCurrentStep] = useState(1);
  const [formData, setFormData] = useState(initialFormData);

  const updateField = (field: string, value: string) => {
    setFormData((prev) => {
      const updated = { ...prev, [field]: value };

      // Auto-configure server settings when email changes
      if (field === "email" && value.includes("@")) {
        const domain = value.split("@")[1]?.toLowerCase();
        const provider = emailProviders[domain];

        if (provider) {
          return {
            ...updated,
            imapServer: provider.imapServer,
            imapPort: provider.imapPort,
            imapSecurity: provider.imapSecurity,
            imapAuth: provider.imapAuth,
            smtpServer: provider.smtpServer,
            smtpPort: provider.smtpPort,
            smtpSecurity: provider.smtpSecurity,
            smtpAuth: provider.smtpAuth,
          };
        }
      }
      return updated;
    });
  };

  const handleNext = () => {
    if (currentStep < 3) {
      setCurrentStep(currentStep + 1);
    }
  };

  const handleBack = () => {
    if (currentStep > 1) {
      setCurrentStep(currentStep - 1);
    }
  };

  const handleSubmit = () => {
    console.log("Account setup complete:", formData);
    alert("Account setup complete! Check console for data.");
    G_AddAccount({
      id: 0, // Will be set by backend
      email: formData.email,
      name: formData.fullName,
      password: formData.password,
      remember_password: formData.rememberPassword === "true",

      imap_username: formData.email,
      imap_password: formData.password,
      imap_host: formData.imapServer,
      imap_port: parseInt(formData.imapPort, 10),
      imap_security: formData.imapSecurity,
      imap_auth_type: formData.imapAuth,
      
      smtp_username: formData.email,
      smtp_password: formData.password,
      smtp_host: formData.smtpServer,
      smtp_port: parseInt(formData.smtpPort, 10),
      smtp_security: formData.smtpSecurity,
      smtp_auth_type: formData.smtpAuth,
      created_at: new Date().toISOString(), // Placeholder,
    }).then(() => {
      console.log("Account added successfully");
    });
  };

  return (
    <div className="flex flex-col gap-6 w-full max-w-125 mx-auto p-4">
      <Card>
        <CardHeader>
          <CardTitle className="text-center">Welcome to Bunny Mail!</CardTitle>
          <CardDescription className="text-center">
            {currentStep === 1 &&
              "Set up your existing email account to get started."}
            {currentStep === 2 &&
              "Configure your imap and smtp server settings."}
            {currentStep === 3 && "Review and confirm your account details."}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <StepIndicator
            steps={[1, 2, 3]}
            currentStep={currentStep}
            stepLabels={["Basic Info", "Server Config", "Review"]}
          />
          {currentStep === 1 && (
            <AccountSetupStep1
              formData={formData}
              updateField={updateField}
              handleNext={handleNext}
            />
          )}
          {currentStep === 2 && (
            <AccountSetupStep2
              formData={formData}
              updateField={updateField}
              handleNext={handleNext}
              handleBack={handleBack}
            />
          )}
          {currentStep === 3 && (
            <AccountSetupStep3
              formData={formData}
              handleBack={handleBack}
              handleSubmit={handleSubmit}
            />
          )}
        </CardContent>
      </Card>
    </div>
  );
}
