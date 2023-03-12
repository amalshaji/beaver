interface IStats {
  active_connections: number;
  cpu_used: number;
  memory_used: number;
  connection_status: IActiveConnectionStatus[];
}

interface IActiveConnectionStatus {
  ID: number;
  Active: boolean;
  LastActiveAt: string;
}

interface ITunnelUser {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string;
  Email: string;
  Active: boolean;
  LastActiveAt: string | null;
}

interface IError {
  error: string;
}
