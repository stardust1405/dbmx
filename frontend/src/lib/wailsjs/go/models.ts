export namespace model {
	
	export class Database {
	    ID: string;
	    PostgresConnectionID: number;
	    PostgresConnectionName: string;
	    Name: string;
	    Colour: string;
	    PoolID: string;
	    IsActive: boolean;
	    Tables: string[];
	    Columns: string[];
	
	    static createFrom(source: any = {}) {
	        return new Database(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.PostgresConnectionID = source["PostgresConnectionID"];
	        this.PostgresConnectionName = source["PostgresConnectionName"];
	        this.Name = source["Name"];
	        this.Colour = source["Colour"];
	        this.PoolID = source["PoolID"];
	        this.IsActive = source["IsActive"];
	        this.Tables = source["Tables"];
	        this.Columns = source["Columns"];
	    }
	}
	export class GenericResponse {
	    ok: boolean;
	    data: any[];
	    rowsAffected: number;
	    message: string;
	    columns: string[];
	    rows: any[][];
	
	    static createFrom(source: any = {}) {
	        return new GenericResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.data = source["data"];
	        this.rowsAffected = source["rowsAffected"];
	        this.message = source["message"];
	        this.columns = source["columns"];
	        this.rows = source["rows"];
	    }
	}
	export class PostgresConnection {
	    ID: number;
	    Name: string;
	    Host: string;
	    Port: string;
	    Username: string;
	    Password: string;
	    Database: string;
	    Env: string;
	    Colour: string;
	    IsActive: boolean;
	
	    static createFrom(source: any = {}) {
	        return new PostgresConnection(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Name = source["Name"];
	        this.Host = source["Host"];
	        this.Port = source["Port"];
	        this.Username = source["Username"];
	        this.Password = source["Password"];
	        this.Database = source["Database"];
	        this.Env = source["Env"];
	        this.Colour = source["Colour"];
	        this.IsActive = source["IsActive"];
	    }
	}
	export class Tab {
	    ID: number;
	    Name: string;
	    Editor: string;
	    Output: string;
	    IsActive: boolean;
	    ActiveDBID?: number;
	    ActiveDB?: string;
	
	    static createFrom(source: any = {}) {
	        return new Tab(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Name = source["Name"];
	        this.Editor = source["Editor"];
	        this.Output = source["Output"];
	        this.IsActive = source["IsActive"];
	        this.ActiveDBID = source["ActiveDBID"];
	        this.ActiveDB = source["ActiveDB"];
	    }
	}

}

