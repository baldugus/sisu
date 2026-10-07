/**
 * Fixture data for the mock backend.
 * Shapes mirror the Go structs in sisu/types/.
 *
 * Scenario: two courses (morning 4 seats, evening 2 seats). Call 1 is closed;
 * call 2 is open with one promotion offer and one waitlist call.
 */

const NAMES = [
  'Ana Paula Ferreira',
  'Bruno Souza Lima',
  'Carla Mendes Santos',
  'Diego Oliveira Costa',
  'Eduarda Rocha Alves',
  'Felipe Nunes Prado',
  'Gabriela Teixeira',
  'Henrique Barros',
  'Isabela Moura',
  'João Victor Ramos',
  'Karen Duarte',
];

export const mockCandidates = NAMES.map((name, i) => ({
  ID: i + 1,
  CPF: `${String(i + 1).padStart(3, '0')}45678900`,
  Name: name,
  SocialName: '',
  BirthDate: '2000-01-01',
  Sex: i % 2 ? 'M' : 'F',
  MotherName: 'Maria da Silva',
  AddressLine: 'Rua das Flores',
  AddressLine2: '',
  HouseNumber: String(10 + i),
  Neighborhood: 'Centro',
  Municipality: 'Rio de Janeiro',
  State: 'RJ',
  CEP: '20000-000',
  Email: `${name.split(' ')[0].toLowerCase()}@email.com`,
  Phone1: `(21) 99999-00${String(i + 1).padStart(2, '0')}`,
  Phone2: '',
}));

export const mockCourses = [
  { ID: 1, Seats: 4, MinimumScore: '650,00', Period: 'morning', Quota: 'Ampla concorrência' },
  { ID: 2, Seats: 2, MinimumScore: '600,00', Period: 'evening', Quota: 'Ampla concorrência' },
];

// [registration ID, course ID, ranking, selection]
const REGS: [number, number, number, 'approved' | 'waitlist'][] = [
  [1, 1, 1, 'approved'],
  [2, 1, 2, 'approved'],
  [3, 1, 3, 'approved'],
  [4, 1, 4, 'approved'],
  [5, 2, 1, 'approved'],
  [6, 2, 2, 'approved'],
  [7, 1, 5, 'waitlist'],
  [8, 1, 6, 'waitlist'],
  [9, 2, 3, 'waitlist'],
  [10, 2, 4, 'waitlist'],
  [11, 1, 7, 'waitlist'],
];

export const mockRegistrations = REGS.map(([id, courseID, ranking, selection]) => ({
  ID: id,
  CourseID: courseID,
  SelectionKind: selection,
  EnrollmentID: `20250000${String(id).padStart(4, '0')}`,
  Option: 1,
  LanguagesScore: '650,00',
  HumanitiesScore: '620,50',
  NaturalSciencesScore: '610,00',
  MathematicsScore: '700,00',
  EssayScore: '880,00',
  CompositeScore: `${700 - ranking * 5},00`,
  Ranking: ranking,
  Candidate: mockCandidates[id - 1],
}));

export const mockCalls = [
  { ID: 1, Number: 1, Status: 'done' },
  { ID: 2, Number: 2, Status: 'calling' },
];

export const mockEntries = [
  // Call 1 — approved list split across semesters.
  { CallID: 1, RegistrationID: 1, Kind: 'initial', Semester: 1, Outcome: 'enrolled', WantsPromotion: false },
  { CallID: 1, RegistrationID: 2, Kind: 'initial', Semester: 1, Outcome: 'absent', WantsPromotion: false },
  { CallID: 1, RegistrationID: 3, Kind: 'initial', Semester: 2, Outcome: 'enrolled', WantsPromotion: true },
  { CallID: 1, RegistrationID: 4, Kind: 'initial', Semester: 2, Outcome: 'enrolled', WantsPromotion: false },
  { CallID: 1, RegistrationID: 5, Kind: 'initial', Semester: 1, Outcome: 'absent', WantsPromotion: false },
  { CallID: 1, RegistrationID: 6, Kind: 'initial', Semester: 2, Outcome: 'enrolled', WantsPromotion: false },
  // Call 2 — promotion first, then the waitlist.
  { CallID: 2, RegistrationID: 3, Kind: 'promotion', Semester: 1, Outcome: 'pending', WantsPromotion: false },
  { CallID: 2, RegistrationID: 9, Kind: 'waitlist', Semester: 1, Outcome: 'pending', WantsPromotion: false },
];

export const mockSemesters = [{ Number: 1 }, { Number: 2 }];

export const approvedSelection = {
  ID: 1,
  Name: 'SISU 2025 — Convocados',
  Kind: 'approved',
  Year: 2025,
  Institution: 'FAETERJ-Rio',
  Degree: 'Tecnológico',
};

export const waitlistSelection = {
  ID: 2,
  Name: 'SISU 2025 — Lista de Espera',
  Kind: 'waitlist',
  Year: 2025,
  Institution: 'FAETERJ-Rio',
  Degree: 'Tecnológico',
};
