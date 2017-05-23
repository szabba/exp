#lang rackjure

(require srfi/1)
(require "termprog.rkt")
(require "ui.rkt")
(require "world-state.rkt")

(define world
  (world-state
   {'jokes 3 'adventures 7}
   (set
    (option "Boast of your deeds in the wild suburbia"
            (λ (qs) (~> (qs 'adventures #:else 0) (>= 3)))

            '("Everyone's impressed you managed to survive amongst the middle class")
            (list (λ (qs) (dict-update qs 'cool-factor #λ(+ % 3) 0))))
    
    (option "Charm them with your sense of humour"
            (λ (qs) (~> (qs 'jokes #:else 0) (>= 2)))

            '("They laugh so hard their stomaches hurt")
            (list (λ (qs) (dict-update qs 'cool-factor #λ(+ % 1) 0))))
    
    (option "Educate the philsitne rabble"
            (λ (qs) (and (~> (qs 'schoolwork #:else 0) (>= 3))
                         (~> (qs 'snobbery #:else 0) (>= 3))))

            '("You'll make some enemies...")
            (list (λ (qs) (dict-update qs 'cool-factor #λ(- % 3) 0))
                  (λ (qs) (dict-update qs 'snobbery #λ(+ % 10) 0)))))))

(run-world! world)