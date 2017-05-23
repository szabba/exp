#lang rackjure

(require srfi/1)

(provide (contract-out
          [struct world-state ((qualities (listof pair?))
                               (options (set/c option?)))]

          [struct option ((title string?)
                          (preconds (listof any/c))
                          (text (listof string?))
                          (effects (listof any/c)))]

          [world-state-available-options (-> world-state? list?)]
          [world-state-apply-option (-> world-state? option? world-state?)]))

(struct world-state (qualities options) #:transparent)
(struct option (title preconds text effects) #:transparent)
(struct precond (qname constraints) #:transparent)

;;; What can be done?

(define (world-state-available-options world)
  (let [(aq (~>> world world-state-qualities))]
    (~>> world
         world-state-options
         set->list
         (filter (λ (opt) (option-available aq opt))))))

(define (option-available aq opt)
  (~>> opt
       option-preconds
       dict->list
       (map pair->precond)
       (every (λ (precond) (eval-precond aq precond)))))

(define (pair->precond pspec)
  (match pspec
    [`(,quality . ,conds) (precond quality conds)]))

(define (eval-precond aq prec)
  (match prec
    [(precond qname constraints)
     (~>> (aq qname #:else 0)
          (eval-constraints constraints))]))

(define (eval-constraints cs v)
  (~>> cs
       (every (λ (c) (eval-constraint c v)))))

(define (eval-constraint c v)
  (match c
    [`(< ,n) (< v n)]
    [`(> ,n) (> v n)]
    [`(<= ,n) (<= v n)]
    [`(>= ,n) (>= v n)]
    [`(= ,n) (= v n)]))

;;; How are things done?

(define (world-state-apply-option world opt)
  ;; Let's just **assume** the option's valid at this point.
  (struct-copy world-state world
               (qualities (~>> world
                               world-state-qualities
                               (apply-effects (~>> opt option-effects))))))

(define (apply-effects effs qs)
  (~>> (fold apply-effect qs effs)
       (filter (compose positive? cdr))))

(define (apply-effect eff qs)
  (match eff
    [(list qname eff-op rest ...)
     (~> eff-op eff-handlers (apply qs qname rest))]))

(define eff-handlers
  {'- (λ (qs qn n)
        (let [(v (qs qn #:else 0))]
          (if (< v n)
              (dict-remove qs qn)
              (qs qn (- v n)))))
      
   '+ (λ (qs qn n)
        (~>> (qs qn #:else 0) (+ n) (qs qn)))
   
   '= (λ (qs qn n) (qs qn n))
   
   'rm (λ (qs qn) (dict-remove qs qn))})
  