<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class TriggerOccurrence
{
    /**
     *
     * @var \DateTime $date
     */
    #[\JMS\Serializer\Annotation\SerializedName('date')]
    public \DateTime $date;

    /**
     *
     * @var ?string $error
     */
    #[\JMS\Serializer\Annotation\SerializedName('error')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?string $error = null;

    /**
     * $event
     *
     * @var array<string, mixed> $event
     */
    #[\JMS\Serializer\Annotation\SerializedName('event')]
    #[\JMS\Serializer\Annotation\Type('array<string, mixed>')]
    public array $event;

    /**
     *
     * @var string $triggerID
     */
    #[\JMS\Serializer\Annotation\SerializedName('triggerID')]
    public string $triggerID;

    /**
     *
     * @var ?WorkflowInstance $workflowInstance
     */
    #[\JMS\Serializer\Annotation\SerializedName('workflowInstance')]
    #[\JMS\Serializer\Annotation\Type('\formance\stack\Models\Shared\WorkflowInstance')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?WorkflowInstance $workflowInstance = null;

    /**
     *
     * @var ?string $workflowInstanceID
     */
    #[\JMS\Serializer\Annotation\SerializedName('workflowInstanceID')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?string $workflowInstanceID = null;

    /**
     * @param  ?\DateTime  $date
     * @param  ?array<string, mixed>  $event
     * @param  ?string  $triggerID
     * @param  ?string  $error
     * @param  ?WorkflowInstance  $workflowInstance
     * @param  ?string  $workflowInstanceID
     */
    public function __construct(?\DateTime $date = null, ?array $event = null, ?string $triggerID = null, ?string $error = null, ?WorkflowInstance $workflowInstance = null, ?string $workflowInstanceID = null)
    {
        $this->date = $date;
        $this->event = $event;
        $this->triggerID = $triggerID;
        $this->error = $error;
        $this->workflowInstance = $workflowInstance;
        $this->workflowInstanceID = $workflowInstanceID;
    }
}