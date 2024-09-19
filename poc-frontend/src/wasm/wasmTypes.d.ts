import { Doxxier, DoxxierPart } from "../types/Doxxier";

declare global {
    export interface Window {
        Go: any;
        CreateDoxxier: () => string;
        UpdateDoxxier: (doxxier:Doxxier) => string;
        AddPart: (part:DoxxierPart) => string;
        RemovePart: (partId:string) => string;
        GetPart: (partId:string) => DoxxierPart;
    }
}

export{}