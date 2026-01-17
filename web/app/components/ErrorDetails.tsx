export interface ErrorDetailsProps {
    error: string;
}

export default function ErrorDetails({ error }: ErrorDetailsProps) {
    return (
        <div>
            <h1>Oops!</h1>
            <p>{error}</p>
        </div>
    );
}