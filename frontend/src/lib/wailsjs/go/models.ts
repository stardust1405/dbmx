export namespace model {
	
	export class Cell {
	    column: string;
	    value: string;
	
	    static createFrom(source: any = {}) {
	        return new Cell(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.column = source["column"];
	        this.value = source["value"];
	    }
	}
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
	export class Output {
	    columns: string[];
	    rows: Cell[][];
	
	    static createFrom(source: any = {}) {
	        return new Output(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.columns = source["columns"];
	        this.rows = this.convertValues(source["rows"], Cell);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
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
	export class QueryResult {
	    ok: boolean;
	    columns: string[];
	    rows: Cell[][];
	    rowsAffected: number;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new QueryResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.columns = source["columns"];
	        this.rows = this.convertValues(source["rows"], Cell);
	        this.rowsAffected = source["rowsAffected"];
	        this.message = source["message"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Tab {
	    ID: number;
	    Name: string;
	    Editor: string;
	    Output: string;
	    IsActive: boolean;
	    ActiveDBID?: string;
	    ActiveDB?: string;
	    ActiveDBColor?: string;
	    Type: string;
	    columns: string[];
	    rows: Cell[][];
	    PostgresConnID?: number;
	    PostgresConnName: string;
	    DBName?: string;
	    Select: string;
	    Limit: string;
	    Offset: string;
	    Where: string;
	    OrderBy: string;
	    GroupBy: string;
	    TableColumns: string;
	    TableColumnsList: string[];
	
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
	        this.ActiveDBColor = source["ActiveDBColor"];
	        this.Type = source["Type"];
	        this.columns = source["columns"];
	        this.rows = this.convertValues(source["rows"], Cell);
	        this.PostgresConnID = source["PostgresConnID"];
	        this.PostgresConnName = source["PostgresConnName"];
	        this.DBName = source["DBName"];
	        this.Select = source["Select"];
	        this.Limit = source["Limit"];
	        this.Offset = source["Offset"];
	        this.Where = source["Where"];
	        this.OrderBy = source["OrderBy"];
	        this.GroupBy = source["GroupBy"];
	        this.TableColumns = source["TableColumns"];
	        this.TableColumnsList = source["TableColumnsList"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

