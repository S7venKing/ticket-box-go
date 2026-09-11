# Ticket Service

This module owns the organizer and event business domain in the Ticket Box Go project.

## Current scope

- Organizer aggregate
- Event aggregate
- Domain validation rules
- MySQL repositories and gRPC server
- Application handlers for organizer/event commands and queries
- Event approval workflow: `draft -> pending_approval -> published`
- Docker migration for the `ticket_db` schema
- Domain tests

## Authorization

Authorization is enforced at the API gateway using the role in the JWT:

- `admin`: creates organizers and approves events
- `organizer`: creates and submits events belonging to its organizer profile
- `user`: cannot manage organizers or events

The ticket service is internal and is called through gRPC by the gateway.
