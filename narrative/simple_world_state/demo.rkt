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
            {'adventures '((>= 3))}
            '("Everyone's impressed you managed to survive amongst the middle class")
            '((cool-factor + 3)))
    (option "Charm them with your sense of humour"
            {'jokes '((>= 2))}
            '("They laugh so hard their stomaches hurt")
            '((cool-factor + 1)))
    (option "Educate the philsitne rabble"
            {'schoolwork '((>= 4)) 'snobbery '((>= 6))}
            '("You'll make some enemies...")
            '((snobbery + 10) (cool-factor - 3))))))

(run-world! world)