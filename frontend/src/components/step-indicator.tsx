import { cn } from "@/lib/utils";
import { Check } from "lucide-react";

export default function StepIndicator({
  steps,
  currentStep,
  stepLabels,
}: {
  steps: number[];
  currentStep: number;
  stepLabels: string[];
}) {
  return (
    <div className="flex items-center justify-center gap-2 mb-6">
      {steps.map((step, index) => (
        <div key={step} className="flex items-center">
          <div className="flex flex-col items-center">
            <div
              className={cn(
                "w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium transition-colors",
                currentStep === step
                  ? "bg-gray-200 text-gray-600"
                  : currentStep > step
                  ? "bg-green-600 text-white"
                  : "bg-accent text-white"
              )}
            >
              {currentStep > step ? <Check className="w-4 h-4" /> : index + 1}
            </div>
            <span className="text-xs mt-1 text-accent-foreground">
              {stepLabels[index]}
            </span>
          </div>
          {index < steps.length - 1 && (
            <div
              className={cn(
                "w-12 h-1 mx-1 mb-5",
                currentStep > step ? "bg-green-600" : "bg-gray-200"
              )}
            />
          )}
        </div>
      ))}
    </div>
  );
}
