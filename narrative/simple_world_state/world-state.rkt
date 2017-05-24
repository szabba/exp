#lang rackjure

(require srfi/1)

(provide (contract-out
          [struct world-state ((qualities any/c)
                               (choices (listof choice?)))]

          [struct choice ((title string?)
                          (preconds (-> any/c boolean?))
                          (text (listof string?))
                          (effects (listof (-> any/c any/c))))]

          [world-state-available-choices (-> world-state? list?)]
          [world-state-apply-choice (-> world-state? choice? world-state?)]))

(struct world-state (qualities choices) #:transparent)
(struct choice (title preconds text effects) #:transparent)

;;; What can be done?

(define (world-state-available-choices world)
  (match world
    [(world-state qualities choices)
     (~>> choices
          (filter (λ (ch)
                    (choice-available? qualities ch))))]))

(define (choice-available? aq opt)
  ((~> opt choice-preconds) aq))

;;; How are things done?

(define (world-state-apply-choice world opt)
  (struct-copy world-state world
               (qualities
                (fold (λ (f x) (f x))
                      (~> world world-state-qualities)
                      (~> opt choice-effects)))))