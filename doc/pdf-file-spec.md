# Generated PDF Specification

This document describes the PDF reports produced by the application and the information exposed in each one.

## Common conventions

The application produces four kinds of PDF:

- website publication list;
- official call/enrollment list;
- email contact list; and
- teacher attendance list.

Each PDF covers one course shift, either `MATUTINO` or `NOTURNO`.

Candidates are separated into sections by competition modality. Modalities with no candidates in the report are omitted. Within each modality, candidates are ordered by their SiSU classification in ascending order, and the displayed sequence starts again at 1.

Candidate scores are shown with two decimal places and a comma as the decimal separator.

### Call report header

The website, enrollment, and email reports share a centered header:

```text
<institution>
<call description> - <year> - <semester>o. Semestre
<course>
TURNO: <shift>
```

For the initial call, the call description is:

```text
CHAMADA REGULAR
```

For subsequent calls, it is:

```text
LISTA DE ESPERA <number>
```

The first waitlist call is numbered 1, corresponding to application call 2. Application call 3 is shown as `LISTA DE ESPERA 2`, and so on.

The header is repeated on every page.

## Website publication list

The website publication list is the public list of candidates included in a particular call and shift.

It is generated in portrait orientation.

For each competition modality, it contains:

| Column | Contents |
| --- | --- |
| `ORDEM` | Sequential position in the displayed modality section. |
| `NOME DO CANDIDATO` | Candidate's legal name. |
| `NOTA` | Candidate's final/composite score. |
| `INSCRIÇÃO DO ENEM` | Candidate's ENEM registration identifier. |

Only candidates assigned to the selected call and shift appear in this report. Their current enrollment or absence status does not remove them from the call's list.

The default filename proposed by the application follows:

```text
website_chamada<call>_<shift>_<year>.pdf
```

Example:

```text
website_chamada1_matutino_2025.pdf
```

## Official call/enrollment list

The official call/enrollment list contains the candidates included in a particular call and provides a space for each candidate's signature.

It is generated in landscape orientation.

For each competition modality, it contains:

| Column | Contents |
| --- | --- |
| `ORDEM` | Sequential position in the displayed modality section. |
| `NOME DO CANDIDATO` | Candidate's legal name. |
| `NOTA` | Candidate's final/composite score. |
| `INSCRIÇÃO DO ENEM` | Candidate's ENEM registration identifier. |
| `ASSINATURA` | Blank line for the candidate's signature. |

Only candidates assigned to the selected call and shift appear in this report. Their current enrollment or absence status does not remove them from the call's list.

The default filename proposed by the application follows:

```text
enrollment_chamada<call>_<shift>_<year>.pdf
```

Example:

```text
enrollment_chamada1_noturno_2025.pdf
```

## Email contact list

The email contact list contains contact information for the candidates included in a particular call.

It is generated in landscape orientation.

For each competition modality, it contains:

| Column | Contents |
| --- | --- |
| `ORDEM` | Sequential position in the displayed modality section. |
| `NOME DO CANDIDATO` | Candidate's legal name. |
| `NOTA` | Candidate's final/composite score. |
| `INSCRIÇÃO DO ENEM` | Candidate's ENEM registration identifier. |
| `EMAIL` | Candidate's email address. |

Only candidates assigned to the selected call and shift appear in this report. Their current enrollment or absence status does not remove them from the call's list.

The default filename proposed by the application follows:

```text
email_chamada<call>_<shift>_<year>.pdf
```

Example:

```text
email_chamada2_matutino_2025.pdf
```

## Teacher attendance list

The teacher attendance list is a general attendance sheet containing all candidates currently marked as enrolled in the selected semester and shift. It is not limited to a particular call.

It is generated in landscape orientation.

Its header is:

```text
LISTA DE PRESENÇA - <institution> <year>.<semester>
TURNO: <shift>
DISCIPLINA:
PROFESSOR:
DATA: ____/____/<year>
```

The semester in the header is the academic semester represented by the report, either `1` or `2`.

For each competition modality, the report contains:

| Column | Contents |
| --- | --- |
| `Seq.` | Sequential position in the displayed modality section. |
| `Nome` | Candidate's legal name. |
| `Assinatura` | Blank line for the candidate's signature. |

The list includes enrolled candidates from every call and from both the approved and waitlist imports, provided they belong to the selected semester and shift.

The default filename proposed by the application follows:

```text
professor_<shift>_<year>.pdf
```

Example:

```text
professor_matutino_2025.pdf
```

## Empty reports

If no candidate matches the selected report and shift, the generated PDF contains its header but no competition-modality sections or candidate rows.
