export interface UserExerciseQuizData {
    [key: string]: string;
}

export interface UserExerciseCodeData {
    name: string;
    content?: string;
    children?: UserExerciseCodeData[];
}