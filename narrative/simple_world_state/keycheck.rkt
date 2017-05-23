#lang rackjure

(require "termprog.rkt")

(define (update key _) key)

(define (view w h model)
  (~>> model ~v list))

(run! (program (void) update view (λ~> (eq? 'escape))))