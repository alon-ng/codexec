import { Progress as ShadcnProgress } from "~/components/ui/progress";
import { cn } from "~/lib/utils";

export interface ProgressProps extends React.ComponentProps<typeof ShadcnProgress> {
    value: number;
    showPercentage?: boolean;
}

function Progress(props: ProgressProps) {
    return (
        <div className="flex flex-col items-end w-full">
            {props.showPercentage && <div className="text-sm text-muted-foreground">{props.value.toFixed(0)}%</div>}
            <ShadcnProgress
                className={cn(props.className, "*:data-[slot=progress-indicator]:bg-linear-to-r *:data-[slot=progress-indicator]:from-codim-purple *:data-[slot=progress-indicator]:to-codim-pink [&>div]:bg-purple-500/20")}
                {...props}
            />
        </div>
    );
}

export { Progress };