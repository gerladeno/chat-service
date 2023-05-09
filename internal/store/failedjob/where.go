// Code generated by ent, DO NOT EDIT.

package failedjob

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/gerladeno/chat-service/internal/store/predicate"
	"github.com/gerladeno/chat-service/internal/types"
)

// ID filters vertices based on their ID field.
func ID(id types.FailedJobID) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id types.FailedJobID) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id types.FailedJobID) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...types.FailedJobID) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...types.FailedJobID) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id types.FailedJobID) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id types.FailedJobID) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id types.FailedJobID) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id types.FailedJobID) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldLTE(FieldID, id))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldEQ(FieldName, v))
}

// Payload applies equality check predicate on the "payload" field. It's identical to PayloadEQ.
func Payload(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldEQ(FieldPayload, v))
}

// Reason applies equality check predicate on the "reason" field. It's identical to ReasonEQ.
func Reason(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldEQ(FieldReason, v))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldEQ(FieldCreatedAt, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldContainsFold(FieldName, v))
}

// PayloadEQ applies the EQ predicate on the "payload" field.
func PayloadEQ(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldEQ(FieldPayload, v))
}

// PayloadNEQ applies the NEQ predicate on the "payload" field.
func PayloadNEQ(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldNEQ(FieldPayload, v))
}

// PayloadIn applies the In predicate on the "payload" field.
func PayloadIn(vs ...string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldIn(FieldPayload, vs...))
}

// PayloadNotIn applies the NotIn predicate on the "payload" field.
func PayloadNotIn(vs ...string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldNotIn(FieldPayload, vs...))
}

// PayloadGT applies the GT predicate on the "payload" field.
func PayloadGT(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldGT(FieldPayload, v))
}

// PayloadGTE applies the GTE predicate on the "payload" field.
func PayloadGTE(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldGTE(FieldPayload, v))
}

// PayloadLT applies the LT predicate on the "payload" field.
func PayloadLT(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldLT(FieldPayload, v))
}

// PayloadLTE applies the LTE predicate on the "payload" field.
func PayloadLTE(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldLTE(FieldPayload, v))
}

// PayloadContains applies the Contains predicate on the "payload" field.
func PayloadContains(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldContains(FieldPayload, v))
}

// PayloadHasPrefix applies the HasPrefix predicate on the "payload" field.
func PayloadHasPrefix(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldHasPrefix(FieldPayload, v))
}

// PayloadHasSuffix applies the HasSuffix predicate on the "payload" field.
func PayloadHasSuffix(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldHasSuffix(FieldPayload, v))
}

// PayloadEqualFold applies the EqualFold predicate on the "payload" field.
func PayloadEqualFold(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldEqualFold(FieldPayload, v))
}

// PayloadContainsFold applies the ContainsFold predicate on the "payload" field.
func PayloadContainsFold(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldContainsFold(FieldPayload, v))
}

// ReasonEQ applies the EQ predicate on the "reason" field.
func ReasonEQ(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldEQ(FieldReason, v))
}

// ReasonNEQ applies the NEQ predicate on the "reason" field.
func ReasonNEQ(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldNEQ(FieldReason, v))
}

// ReasonIn applies the In predicate on the "reason" field.
func ReasonIn(vs ...string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldIn(FieldReason, vs...))
}

// ReasonNotIn applies the NotIn predicate on the "reason" field.
func ReasonNotIn(vs ...string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldNotIn(FieldReason, vs...))
}

// ReasonGT applies the GT predicate on the "reason" field.
func ReasonGT(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldGT(FieldReason, v))
}

// ReasonGTE applies the GTE predicate on the "reason" field.
func ReasonGTE(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldGTE(FieldReason, v))
}

// ReasonLT applies the LT predicate on the "reason" field.
func ReasonLT(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldLT(FieldReason, v))
}

// ReasonLTE applies the LTE predicate on the "reason" field.
func ReasonLTE(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldLTE(FieldReason, v))
}

// ReasonContains applies the Contains predicate on the "reason" field.
func ReasonContains(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldContains(FieldReason, v))
}

// ReasonHasPrefix applies the HasPrefix predicate on the "reason" field.
func ReasonHasPrefix(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldHasPrefix(FieldReason, v))
}

// ReasonHasSuffix applies the HasSuffix predicate on the "reason" field.
func ReasonHasSuffix(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldHasSuffix(FieldReason, v))
}

// ReasonEqualFold applies the EqualFold predicate on the "reason" field.
func ReasonEqualFold(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldEqualFold(FieldReason, v))
}

// ReasonContainsFold applies the ContainsFold predicate on the "reason" field.
func ReasonContainsFold(v string) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldContainsFold(FieldReason, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.FailedJob {
	return predicate.FailedJob(sql.FieldLTE(FieldCreatedAt, v))
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.FailedJob) predicate.FailedJob {
	return predicate.FailedJob(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.FailedJob) predicate.FailedJob {
	return predicate.FailedJob(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.FailedJob) predicate.FailedJob {
	return predicate.FailedJob(func(s *sql.Selector) {
		p(s.Not())
	})
}
