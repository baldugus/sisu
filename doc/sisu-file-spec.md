# SiSU File Specification

This document describes the structure and contents of the files supplied by MEC through SiSU that are relevant to this application.

## Approved Candidates File

The approved candidates file contains the candidates selected in the regular SiSU result for an institution. Each row represents one candidate's registration for a course, shift, and competition modality.

The file covers the annual selection. Its seat counts and rankings are therefore annual rather than specific to one academic semester.

### Format

Files received by the institution have the following characteristics:

- Text is encoded as UTF-8 and may begin with a byte order mark (BOM).
- The header fields are separated by commas.
- The values in each data row are separated by semicolons.
- Records may be separated by carriage returns (`CR`, `\r`).
- Values containing semicolons may be surrounded by double quotes.

### Fields used by this application

The SiSU file contains more fields than the application requires. The fields below form the subset of the external file contract used by SISU.

#### Institution and course

| Field | Contents |
| --- | --- |
| `SG_IES` | Abbreviation of the higher-education institution. |
| `NO_CURSO` | Name of the course or degree program. |
| `DS_TURNO` | Shift in which the course is offered. The observed values are `Matutino` and `Noturno`. |

#### Competition modality and seats

| Field | Contents |
| --- | --- |
| `NO_MODALIDADE_CONCORRENCIA` | Name or description of the competition modality, such as broad competition or a quota category. |
| `QT_VAGAS_CONCORRENCIA` | Number of annual seats SiSU allocated to the course, shift, and competition modality represented by the row. |
| `NU_NOTACORTE_CONCORRIDA` | Cutoff score for that competition group. |
| `NU_CLASSIFICACAO` | Candidate's placement in that competition group. Lower positive numbers represent better placement. |

Candidates competing for the same course, shift, and competition modality are expected to share the same seat count and cutoff score. Their classifications identify their relative order within that group.

#### Registration

| Field | Contents |
| --- | --- |
| `CO_INSCRICAO_ENEM` | Candidate's ENEM registration identifier. |
| `ST_OPCAO` | Number of the course option selected by the candidate. |

#### ENEM scores

| Field | Contents |
| --- | --- |
| `NU_NOTA_L` | Languages score. |
| `NU_NOTA_CH` | Humanities score. |
| `NU_NOTA_CN` | Natural sciences score. |
| `NU_NOTA_M` | Mathematics score. |
| `NU_NOTA_R` | Essay score. |
| `NU_NOTA_CANDIDATO` | Candidate's final or composite score for the course. |

Score fields are decimal numbers and may use either a comma or a period as the decimal separator.

#### Candidate identity

| Field | Contents |
| --- | --- |
| `NU_CPF_INSCRITO` | Candidate's CPF. |
| `NO_INSCRITO` | Candidate's legal name. |
| `NO_SOCIAL` | Candidate's social name, when provided. |
| `DT_NASCIMENTO` | Candidate's date of birth. |
| `TP_SEXO` | Candidate's sex. |
| `NO_MAE` | Candidate's mother's name. |

#### Address and contact information

| Field | Contents |
| --- | --- |
| `DS_LOGRADOURO` | Street or other primary address line. |
| `NU_ENDERECO` | Address number. |
| `DS_COMPLEMENTO` | Address complement. |
| `NO_BAIRRO` | Neighborhood. |
| `NO_MUNICIPIO` | Municipality. |
| `SG_UF_INSCRITO` | State abbreviation. |
| `NU_CEP` | Postal code. |
| `DS_EMAIL` | Email address. |
| `NU_FONE1` | Primary phone number. |
| `NU_FONE2` | Secondary phone number. |

### Other fields

The approved candidates file also contains administrative, socioeconomic, affirmative-action, enrollment, and course-specific score fields that are not part of the subset used by this application.

The following is the complete header observed in the supplied files:

```text
NU_ETAPA,CO_IES,NO_IES,SG_IES,SG_UF_IES,NO_CAMPUS,CO_IES_CURSO,NO_CURSO,DS_TURNO,DS_FORMACAO,QT_VAGAS_CONCORRENCIA,CO_INSCRICAO_ENEM,NO_INSCRITO,NO_SOCIAL,NU_CPF_INSCRITO,DT_NASCIMENTO,TP_SEXO,NU_RG,NO_MAE,DS_LOGRADOURO,NU_ENDERECO,DS_COMPLEMENTO,SG_UF_INSCRITO,NO_MUNICIPIO,NO_BAIRRO,NU_CEP,NU_FONE1,NU_FONE2,DS_EMAIL,NU_NOTA_L,NU_NOTA_CH,NU_NOTA_CN,NU_NOTA_M,NU_NOTA_R,CO_CURSO_INSCRICAO,ST_OPCAO,NO_MODALIDADE_CONCORRENCIA,ST_BONUS_PERC,QT_BONUS_PERCENTUAL,NO_ACAO_AFIRMATIVA_BONUS,NO_ACAO_AFIRMATIVA_PROPRIA_IES,NU_NOTA_CANDIDATO,NU_NOTACORTE_CONCORRIDA,NU_CLASSIFICACAO,DS_MATRICULA,DT_MATRICULA_EFETIVADA,ENSINO_MEDIO,COR_RACA,QUILOMBOLA,PCD,ST_RANK_ENSINO_MEDIO,ST_RANK_RACA,ST_RANK_QUILOMBOLA,ST_RANK_PcD,ST_CONFIRMA_LGPD,TOTAL_MEMBROS_FAMILIAR,RENDA_FAMILIAR_BRUTA,SALARIO_MINIMO,PERFIL_ECONOMICO_LEI_COTAS,TIPO_CONCORRENCIA,DT_CURSO_INSCRICAO,NU_NOTA_CURSO_L,NU_NOTA_CURSO_CH,NU_NOTA_CURSO_CN,NU_NOTA_CURSO_M,NU_NOTA_CURSO_R,ST_APROVADO,DT_MES_DIA_MATRICULA,ST_MATRICULA_CANCELADA,DT_MATRICULA_CANCELADA,VAGA_REMANEJADA,NO_ACAO_AFIRMATIVA_PROPRIA_IES
```

## Waitlisted Candidates File

We did not have access to this file yet so for now we assume this is similar to the approved file.
