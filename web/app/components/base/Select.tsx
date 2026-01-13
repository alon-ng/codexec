import * as React from "react";
import {
    Select as SelectPrimitive,
    SelectContent,
    SelectGroup,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "~/components/ui/select";
import { cn } from "~/lib/utils";

export interface BaseSelectOption {
    value: string;
    label: string;
    [key: string]: any;
}

interface BaseSelectProps {
    value: string;
    onValueChange: (value: string) => void;
    options: BaseSelectOption[];
    placeholder?: string;
    className?: string;
    triggerClassName?: string;
    contentClassName?: string;
    renderTrigger?: (option: BaseSelectOption | undefined) => React.ReactNode;
    renderOption?: (option: BaseSelectOption) => React.ReactNode;
    position?: "item-aligned" | "popper";
}

export function Select({
    value,
    onValueChange,
    options,
    placeholder,
    className,
    triggerClassName,
    contentClassName,
    renderTrigger,
    renderOption,
    position = "popper",
}: BaseSelectProps) {
    const selectedOption = options.find((opt) => opt.value === value);

    return (
        <SelectPrimitive value={value} onValueChange={onValueChange}>
            <SelectTrigger className={cn(triggerClassName, className)}>
                {renderTrigger ? (
                    renderTrigger(selectedOption)
                ) : (
                    <SelectValue placeholder={placeholder}>
                        {selectedOption?.label}
                    </SelectValue>
                )}
            </SelectTrigger>
            <SelectContent
                className={cn(contentClassName, "select-blur-in")}
                position={position}
            >
                <SelectGroup>
                    {options.map((option) => (
                        <SelectItem key={option.value} value={option.value}>
                            {renderOption ? renderOption(option) : option.label}
                        </SelectItem>
                    ))}
                </SelectGroup>
            </SelectContent>
        </SelectPrimitive>
    );
}
