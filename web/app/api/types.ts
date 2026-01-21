export interface UserExerciseQuizData {
    [key: string]: string;
}

export interface UserExerciseCodeData {
    name: string;
    directories: Directory[];
    files: File[];
}

export interface File {
    name: string;
    ext: string;
    content: string;
}

export interface Directory {
    name: string;
    directories: Directory[];
    files: File[];
}