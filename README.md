# README

## About

This project aims to develop an application to streamline the management of student admissions at the Faculdade de Educação Técnologica do Estado do Rio de Janeiro Fernando Motta institution. The app processes data provided by the Ministry of Education’s Sistema de Seleção Unificada system and provides the means to manage applicants, calls and reports.

The original code was somewhat messy and hastily assembled due to the tight deadline for the Minimum Viable Product (MVP), bearing a resemblance to the [eXtreme Go Horse Methodology](https://medium.com/@dekaah/22-axioms-of-the-extreme-go-horse-methodology-xgh-9fa739ab55b4).

A rewrite is underway and has landed substantial pieces already: yearly (rather than per-semester) imports with automatic semester splitting, a typed Wails API boundary (bound methods return concrete typed values instead of an untyped `Data any` payload), and a frontend UI migration to Tailwind v4 + shadcn/ui. Existing functionality is retained throughout.

## Tech Stack

- **Backend**: Go 1.24+ with the Wails v2 framework
- **Frontend**: React 18 + TypeScript + Vite (Tailwind v4 + shadcn/ui)
- **Database**: SQLite with go-jet for queries and golang-migrate for migrations

See **`CLAUDE.md`** for build/development commands and architecture notes, and **`docs/`** for the testing guide and design backlog.
