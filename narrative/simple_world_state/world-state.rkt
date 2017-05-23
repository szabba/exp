#lang rackjure

(require srfi/1)

(provide (contract-out
          [struct world-state ((qualities any/c)
                               ; TODO/FIXME: listof ?
                               (options (set/c option?)))]

          [struct option ((title string?)
                          (preconds (-> any/c boolean?))
                          (text (listof string?))
                          (effects (listof (-> any/c any/c))))]

          [world-state-available-options (-> world-state? list?)]
          [world-state-apply-option (-> world-state? option? world-state?)]))

(struct world-state (qualities options) #:transparent)
(struct option (title preconds text effects) #:transparent)

;;; What can be done?

(define (world-state-available-options world)
  (let [(aq (~>> world world-state-qualities))]
    (~>> world
         world-state-options
         set->list
         (filter (λ (opt) (option-available aq opt))))))

(define (option-available aq opt)
  ((~> opt option-preconds) aq))
       
;;; How are things done?

(define (world-state-apply-option world opt)
  (struct-copy world-state world
               (qualities
                (fold (λ (f x) (f x))
                      (~> world world-state-qualities)
                      (~> opt option-effects)))))