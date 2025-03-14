export namespace model {
	
	export class GenericResponse {
	    ok: boolean;
	    data: any[];
	    rowsAffected: number;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new GenericResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.data = source["data"];
	        this.rowsAffected = source["rowsAffected"];
	        this.message = source["message"];
	    }
	}
	export class Postgres {
	    ID: number;
	    Name: string;
	    Host: string;
	    Port: string;
	    Username: string;
	    Password: string;
	    Database: string;
	    Env: string;
	    Colour?: string;
	
	    static createFrom(source: any = {}) {
	        return new Postgres(source);
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
	    }
	}

}

