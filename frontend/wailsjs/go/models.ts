export namespace types {
	
	export class Call {
	    ID: number;
	    Status: string;
	    Number: number;
	    SemesterID: number;
	
	    static createFrom(source: any = {}) {
	        return new Call(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Status = source["Status"];
	        this.Number = source["Number"];
	        this.SemesterID = source["SemesterID"];
	    }
	}
	export class Candidate {
	    ID: number;
	    CPF: string;
	    Name: string;
	    SocialName: string;
	    BirthDate: string;
	    Sex: string;
	    MotherName: string;
	    AddressLine: string;
	    AddressLine2: string;
	    HouseNumber: string;
	    Neighborhood: string;
	    Municipality: string;
	    State: string;
	    CEP: string;
	    Email: string;
	    Phone1: string;
	    Phone2: string;
	
	    static createFrom(source: any = {}) {
	        return new Candidate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.CPF = source["CPF"];
	        this.Name = source["Name"];
	        this.SocialName = source["SocialName"];
	        this.BirthDate = source["BirthDate"];
	        this.Sex = source["Sex"];
	        this.MotherName = source["MotherName"];
	        this.AddressLine = source["AddressLine"];
	        this.AddressLine2 = source["AddressLine2"];
	        this.HouseNumber = source["HouseNumber"];
	        this.Neighborhood = source["Neighborhood"];
	        this.Municipality = source["Municipality"];
	        this.State = source["State"];
	        this.CEP = source["CEP"];
	        this.Email = source["Email"];
	        this.Phone1 = source["Phone1"];
	        this.Phone2 = source["Phone2"];
	    }
	}
	export class Course {
	    ID: number;
	    Seats: number;
	    MinimumScore?: string;
	    Period: string;
	    Quota: string;
	
	    static createFrom(source: any = {}) {
	        return new Course(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Seats = source["Seats"];
	        this.MinimumScore = source["MinimumScore"];
	        this.Period = source["Period"];
	        this.Quota = source["Quota"];
	    }
	}
	export class Registration {
	    ID: number;
	    EnrollmentID: string;
	    Option: number;
	    LanguagesScore?: string;
	    HumanitiesScore?: string;
	    NaturalSciencesScore?: string;
	    MathematicsScore?: string;
	    EssayScore?: string;
	    CompositeScore?: string;
	    Ranking: number;
	    Status: string;
	    Candidate?: Candidate;
	    SemesterID?: number;
	
	    static createFrom(source: any = {}) {
	        return new Registration(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.EnrollmentID = source["EnrollmentID"];
	        this.Option = source["Option"];
	        this.LanguagesScore = source["LanguagesScore"];
	        this.HumanitiesScore = source["HumanitiesScore"];
	        this.NaturalSciencesScore = source["NaturalSciencesScore"];
	        this.MathematicsScore = source["MathematicsScore"];
	        this.EssayScore = source["EssayScore"];
	        this.CompositeScore = source["CompositeScore"];
	        this.Ranking = source["Ranking"];
	        this.Status = source["Status"];
	        this.Candidate = this.convertValues(source["Candidate"], Candidate);
	        this.SemesterID = source["SemesterID"];
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
	export class RegistrationDetail {
	    Registration?: Registration;
	    Course?: Course;
	    Call?: Call;
	
	    static createFrom(source: any = {}) {
	        return new RegistrationDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Registration = this.convertValues(source["Registration"], Registration);
	        this.Course = this.convertValues(source["Course"], Course);
	        this.Call = this.convertValues(source["Call"], Call);
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
	export class Score {
	    Value: number;
	
	    static createFrom(source: any = {}) {
	        return new Score(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Value = source["Value"];
	    }
	}
	export class Selection {
	    ID: number;
	    Name: string;
	    Kind: string;
	    Year: number;
	    Institution: string;
	    Degree: string;
	
	    static createFrom(source: any = {}) {
	        return new Selection(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Name = source["Name"];
	        this.Kind = source["Kind"];
	        this.Year = source["Year"];
	        this.Institution = source["Institution"];
	        this.Degree = source["Degree"];
	    }
	}
	export class Semester {
	    ID: number;
	    Year: number;
	    Number: number;
	    Status: number;
	
	    static createFrom(source: any = {}) {
	        return new Semester(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Year = source["Year"];
	        this.Number = source["Number"];
	        this.Status = source["Status"];
	    }
	}

}

