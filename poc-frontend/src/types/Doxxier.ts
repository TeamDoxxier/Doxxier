interface Doxxier {
    Id: string;
    Original: string;
    Packaged: string;
    EstimatedDelivery: string;
    IsRevealed: boolean;
    IsTransferring: boolean;
    Progress: number; // Overall transfer progress
    Description: string;
    Thumbnails: FileThumbnail[];
    ThumbnailProgress: number[]; // Individual file progress
    Parts: Array<DoxxierPart>;
    CreatedAt: Date
  }

  interface DoxxierPart {
    Id: string; 
    Content: ArrayBuffer;
    Context: string; 
  }

  interface FileThumbnail {
    Url: string;
    Name: string;
    IsImage: boolean;
  }
  
  export type { Doxxier, DoxxierPart, FileThumbnail };