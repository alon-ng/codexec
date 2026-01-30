import { Markdown } from "@tiptap/markdown";
import { EditorContent, useEditor } from "@tiptap/react";
import StarterKit from "@tiptap/starter-kit";
import { UserRound } from "lucide-react";
import type { ModelsChatMessage } from "~/api/generated/model";
import codyAvatar from "~/assets/cody-256.png";
import ThreeDots from "~/assets/three-dots.svg?react";
import { Skeleton } from "~/components/ui/skeleton";
import { cn } from "~/lib/utils";
import { prose } from '~/utils/prose';

export interface ChatMessageProps {
    message?: ModelsChatMessage;
    isLoadingCodeResponse?: boolean;
}

export default function ChatMessage({ message, isLoadingCodeResponse }: ChatMessageProps) {
    if (isLoadingCodeResponse) {
        return (
            <div className={cn("flex flex-row gap-2 w-full")}>
                <img src={codyAvatar} className="size-9 bg-white rounded-full p-1 shadow-sm" />
                <span className="text-xs font-medium shadow-sm rounded-lg p-2 bg-background overflow-hidden break-words whitespace-normal min-w-0">
                    <ThreeDots className="size-4" />
                </span>
            </div >
        )
    }

    if (!message) {
        return (
            <div className={cn("flex flex-row gap-2 w-full")}>
                <Skeleton className="size-9rounded-full shadow-sm" />
                <Skeleton className="flex-1 h-4 shadow-sm rounded-lg min-w-0" />
            </div >
        )
    }

    const isCody = message.role === "assistant";

    const editor = useEditor({
        extensions: [StarterKit, Markdown],
        editable: false,
        immediatelyRender: false,
        content: message.content,
        contentType: 'markdown',
    });

    return (
        <div className={cn("flex gap-2 w-full", isCody ? "flex-row" : "flex-row-reverse")}>
            {isCody ? (
                <img src={codyAvatar} className="size-9 bg-white rounded-full p-1 shadow-sm shrink-0" />
            ) : (
                <UserRound className="size-9 bg-white rounded-full p-1 shadow-sm shrink-0" />
            )}
            <div className={cn("text-xs font-medium shadow-sm rounded-lg p-2 bg-background overflow-hidden break-words whitespace-normal min-w-0", prose)}>
                <EditorContent editor={editor} />
            </div>
        </div>
    );
}

